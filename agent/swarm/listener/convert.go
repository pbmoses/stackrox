package listener

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/docker"
	"bitbucket.org/stack-rox/apollo/pkg/images"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type serviceWrap swarm.Service

func (s serviceWrap) asDeployment(client *client.Client) *v1.Deployment {
	var updatedTime *timestamp.Timestamp
	up := s.UpdateStatus
	if up != nil && up.CompletedAt != nil {
		var err error
		updatedTime, err = ptypes.TimestampProto(*up.CompletedAt)
		if err != nil {
			log.Error(err)
		}
	}

	image := images.GenerateImageFromString(s.Spec.TaskTemplate.ContainerSpec.Image)

	retries := 0
	for image.Sha == "" && retries <= 15 {
		time.Sleep(time.Second)
		image.Sha = s.getSHAFromTask(client)
		retries++
	}
	if image.Sha == "" {
		log.Warnf("Couldn't find an image SHA for service %s", s.ID)
	}

	m := modeWrap(s.Spec.Mode)

	return &v1.Deployment{
		Id:        s.ID,
		Name:      s.Spec.Name,
		Version:   fmt.Sprintf("%d", s.Version.Index),
		Type:      m.asType(),
		Replicas:  m.asReplica(),
		UpdatedAt: updatedTime,
		Containers: []*v1.Container{
			{
				Config:          s.getContainerConfig(),
				Image:           image,
				SecurityContext: s.getSecurityContext(),
				Volumes:         s.getVolumes(),
				Ports:           s.getPorts(),
			},
		},
	}
}

func (s serviceWrap) getContainerConfig() *v1.ContainerConfig {
	spec := s.Spec.TaskTemplate.ContainerSpec

	envSlice := make([]*v1.ContainerConfig_EnvironmentConfig, 0, len(spec.Env))
	for _, env := range spec.Env {
		parts := strings.SplitN(env, `=`, 2)
		if len(parts) == 2 {
			envSlice = append(envSlice, &v1.ContainerConfig_EnvironmentConfig{
				Key:   parts[0],
				Value: parts[1],
			})
		}
	}

	return &v1.ContainerConfig{
		Args:      spec.Args,
		Command:   spec.Command,
		Directory: spec.Dir,
		Env:       envSlice,
		User:      spec.User,
	}
}

func (s serviceWrap) getSecurityContext() *v1.SecurityContext {
	spec := s.Spec.TaskTemplate.ContainerSpec

	if spec.Privileges == nil || spec.Privileges.SELinuxContext == nil {
		return nil
	}

	return &v1.SecurityContext{
		Selinux: &v1.SecurityContext_SELinux{
			User:  spec.Privileges.SELinuxContext.User,
			Role:  spec.Privileges.SELinuxContext.Role,
			Type:  spec.Privileges.SELinuxContext.Type,
			Level: spec.Privileges.SELinuxContext.Level,
		},
	}
}

func (s serviceWrap) getPorts() []*v1.PortConfig {
	output := make([]*v1.PortConfig, len(s.Endpoint.Ports))
	for i, p := range s.Endpoint.Ports {
		output[i] = &v1.PortConfig{
			Name:          p.Name,
			ContainerPort: int32(p.TargetPort),
			Protocol:      string(p.Protocol),
		}
	}

	return output
}

func (s serviceWrap) getVolumes() []*v1.Volume {
	spec := s.Spec.TaskTemplate.ContainerSpec

	output := make([]*v1.Volume, len(spec.Mounts))

	for i, m := range spec.Mounts {
		output[i] = &v1.Volume{
			Name:     m.Source,
			Path:     m.Target,
			Type:     string(m.Type),
			ReadOnly: m.ReadOnly,
		}
	}

	for _, secret := range spec.Secrets {
		path := ""
		if secret.File != nil {
			path = `/run/secrets/` + secret.File.Name
		}
		output = append(output, &v1.Volume{
			Name: secret.SecretName,
			Path: path,
			Type: `secret`,
		})
	}
	return output
}

func (s serviceWrap) getSHAFromTask(client *client.Client) string {
	opts := filters.NewArgs()
	opts.Add("service", s.ID)
	opts.Add("desired-state", "running")
	ctx, cancel := docker.TimeoutContext()
	defer cancel()
	tasks, err := client.TaskList(ctx, types.TaskListOptions{Filters: opts})
	if err != nil {
		log.Errorf("Couldn't enumerate service %s tasks to get image SHA: %s", s.ID, err)
		return ""
	}
	for _, t := range tasks {
		id := t.Status.ContainerStatus.ContainerID
		if id == "" {
			continue
		}
		ctx, cancel := docker.TimeoutContext()
		defer cancel()
		container, err := client.ContainerInspect(ctx, id)
		if err != nil {
			log.Warnf("Couldn't inspect %s to get image SHA for service %s: %s", id, s.ID, err)
			continue
		}
		// TODO(cg): If the image is specified only as a tag, and Swarm can't
		// resolve to a SHA256 digest when launching, the image SHA may actually
		// differ between tasks on different nodes.
		return container.Image
	}
	return ""
}

type modeWrap swarm.ServiceMode

func (m modeWrap) asType() string {
	if m.Replicated != nil {
		return `Replicated`
	}

	return `Global`
}

func (m modeWrap) asReplica() int64 {
	if m.Replicated != nil {
		return int64(*m.Replicated.Replicas)
	}

	return 0
}
