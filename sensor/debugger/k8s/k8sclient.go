package k8s

import (
	"context"
	"testing"

	appVersioned "github.com/openshift/client-go/apps/clientset/versioned"
	configVersioned "github.com/openshift/client-go/config/clientset/versioned"
	routeVersioned "github.com/openshift/client-go/route/clientset/versioned"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	k8sConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

// MakeFakeClient creates a k8s client that is not connected to any cluster
func MakeFakeClient() *ClientSet {
	return &ClientSet{
		k8s: fake.NewSimpleClientset(),
	}
}

// ClientSet is a test version of kubernetes.ClientSet
type ClientSet struct {
	dynamic         dynamic.Interface
	k8s             kubernetes.Interface
	openshiftApps   appVersioned.Interface
	openshiftConfig configVersioned.Interface
	openshiftRoute  routeVersioned.Interface
}

// MakeOutOfClusterClient creates a k8s client that uses host configuration to connect to a cluster.
// If host machine has a KUBECONFIG env set it will use it to connect to the respective cluster.
func MakeOutOfClusterClient() (*ClientSet, error) {
	config, err := k8sConfig.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "getting k8s config")
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "creating ClientSet")
	}

	return &ClientSet{
		k8s: k8sClient,
	}, nil
}

// Kubernetes returns the kubernetes interface
func (c *ClientSet) Kubernetes() kubernetes.Interface {
	return c.k8s
}

// OpenshiftApps returns the OpenshiftApps interface
// This is not used in tests!
func (c *ClientSet) OpenshiftApps() appVersioned.Interface {
	return c.openshiftApps
}

// OpenshiftConfig returns the OpenshiftConfig interface
// This is not used in tests!
func (c *ClientSet) OpenshiftConfig() configVersioned.Interface {
	return c.openshiftConfig
}

// OpenshiftRoute returns the OpenshiftRoute interface
// This is not used in tests!
func (c *ClientSet) OpenshiftRoute() routeVersioned.Interface {
	return c.openshiftRoute
}

// Dynamic returns the Dynamic interface
// This is not used in tests!
func (c *ClientSet) Dynamic() dynamic.Interface {
	return c.dynamic
}

// SetupExampleCluster creates a fake node and default namespace in the fake k8s client.
func (c *ClientSet) SetupExampleCluster(t *testing.T) {
	_, err := c.k8s.CoreV1().Nodes().Create(context.Background(), &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "go-test-node",
		},
		Spec:   v1.NodeSpec{},
		Status: v1.NodeStatus{},
	}, metav1.CreateOptions{})

	require.NoError(t, err)

	_, err = c.k8s.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec:   v1.NamespaceSpec{},
		Status: v1.NamespaceStatus{},
	}, metav1.CreateOptions{})

	require.NoError(t, err)
}

// MustCreateRole creates a k8s role. Test fails if creation fails.
func (c *ClientSet) MustCreateRole(t *testing.T, name string) {
	_, err := c.k8s.RbacV1().Roles("default").Create(context.Background(), &v12.Role{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Rules: []v12.PolicyRule{{
			APIGroups: []string{""},
			Resources: []string{""},
			Verbs:     []string{"get"},
		}, {
			APIGroups: []string{""},
			Resources: []string{""},
			Verbs:     []string{"list"},
		}},
	}, metav1.CreateOptions{})
	require.NoError(t, err)
}

// MustCreateRoleBinding creates a k8s role binding. Test fails if creation fails.
func (c *ClientSet) MustCreateRoleBinding(t *testing.T, bindingName, roleName, serviceAccountName string) {
	_, err := c.k8s.RbacV1().RoleBindings("default").Create(context.Background(), &v12.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bindingName,
			Namespace: "default",
		},
		Subjects: []v12.Subject{
			{
				Name:      serviceAccountName,
				Kind:      "ServiceAccount",
				Namespace: "default",
			},
		},
		RoleRef: v12.RoleRef{
			APIGroup: "",
			Kind:     "",
			Name:     roleName,
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)
}

// DeploymentOpts defines a function type for changing deployments before being applied.
type DeploymentOpts func(obj *appsv1.Deployment)

// WithServiceAccountName injects service account name to deployment before applying.
func DeploymentWithServiceAccountName(name string) DeploymentOpts {
	return func(obj *appsv1.Deployment) {
		obj.Spec.Template.Spec.ServiceAccountName = name
	}
}

// MustCreateDeployment creates a k8s deployment. Test fails if creation fails.
func (c *ClientSet) MustCreateDeployment(t *testing.T, name string, opts ...DeploymentOpts) {
	d1 := &appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: appsv1.DeploymentSpec{
			Replicas: nil,
			Selector: nil,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: v1.PodSpec{
					Volumes:                       nil,
					InitContainers:                nil,
					Containers:                    nil,
					EphemeralContainers:           nil,
					RestartPolicy:                 "",
					TerminationGracePeriodSeconds: nil,
					ActiveDeadlineSeconds:         nil,
					DNSPolicy:                     "",
					NodeSelector:                  nil,
					ServiceAccountName:            "",
					DeprecatedServiceAccount:      "",
					AutomountServiceAccountToken:  nil,
					NodeName:                      "",
					HostNetwork:                   false,
					HostPID:                       false,
					HostIPC:                       false,
					ShareProcessNamespace:         nil,
					SecurityContext:               nil,
					ImagePullSecrets:              nil,
					Hostname:                      "",
					Subdomain:                     "",
					Affinity:                      nil,
					SchedulerName:                 "",
					Tolerations:                   nil,
					HostAliases:                   nil,
					PriorityClassName:             "",
					Priority:                      nil,
					DNSConfig:                     nil,
					ReadinessGates:                nil,
					RuntimeClassName:              nil,
					EnableServiceLinks:            nil,
					PreemptionPolicy:              nil,
					Overhead:                      nil,
					TopologySpreadConstraints:     nil,
					SetHostnameAsFQDN:             nil,
					OS:                            nil,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type:          "",
				RollingUpdate: nil,
			},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    nil,
			Paused:                  false,
			ProgressDeadlineSeconds: nil,
		},
		Status: appsv1.DeploymentStatus{
			ObservedGeneration:  0,
			Replicas:            0,
			UpdatedReplicas:     0,
			ReadyReplicas:       0,
			AvailableReplicas:   0,
			UnavailableReplicas: 0,
			Conditions:          nil,
			CollisionCount:      nil,
		},
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nginx",
				},
				Spec: v1.PodSpec{
					Volumes:        nil,
					InitContainers: nil,
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.14.2",
							Ports: []v1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}

	for _, opt := range opts {
		opt(deployment)
	}

	_, err := c.k8s.AppsV1().Deployments("default").Create(context.Background(), deployment, metav1.CreateOptions{})
	require.NoError(t, err)
}

// PodOpts defines a function type for changing deployments before being applied.
type PodOpts func(obj *v1.Pod)

// PodWithServiceAccountName injects service account name to pod before applying.
func PodWithServiceAccountName(name string) PodOpts {
	return func(obj *v1.Pod) {
		obj.Spec.ServiceAccountName = name
	}
}

// PodWithLabels injects Labels to Pod before applying.
func PodWithLabels(labels map[string]string) PodOpts {
	return func(obj *v1.Pod) {
		obj.ObjectMeta.Labels = labels
	}
}

//MustCreatePod creates a k8s pod. Test fails if creation fails
func (c *ClientSet) MustCreatePod(t *testing.T, name string, opts ...PodOpts) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels:    map[string]string{"app": "test1"},
		},
		Spec: v1.PodSpec{
			Volumes:        nil,
			InitContainers: nil,
			Containers: []v1.Container{
				{
					Name:  "nginx",
					Image: "nginx:1.14.2",
					Ports: []v1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}

	for _, opt := range opts {
		opt(pod)
	}

	_, err := c.k8s.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	require.NoError(t, err)
}

// ServiceOpts defines a function type for changing deployments before being applied.
type ServiceOpts func(obj *v1.Service)

// ServiceWithSelector injects Selector to Service before applying.
func ServiceWithSelector(selector map[string]string) ServiceOpts {
	return func(obj *v1.Service) {
		obj.Spec.Selector = selector
	}
}

//MustCreateService creates a k8s Service. Test fails if creation fails
func (c *ClientSet) MustCreateService(t *testing.T, name string, opts ...ServiceOpts) {
	service := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       name,
			GenerateName:               "",
			Namespace:                  "default",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ClusterName:                "",
			ManagedFields:              nil,
		},
		Spec: v1.ServiceSpec{
			Ports:                         nil,
			Selector:                      map[string]string{"app": "test1"},
			ClusterIP:                     "",
			ClusterIPs:                    nil,
			Type:                          "",
			ExternalIPs:                   nil,
			SessionAffinity:               "",
			LoadBalancerIP:                "",
			LoadBalancerSourceRanges:      nil,
			ExternalName:                  "",
			ExternalTrafficPolicy:         "",
			HealthCheckNodePort:           0,
			PublishNotReadyAddresses:      false,
			SessionAffinityConfig:         nil,
			IPFamilies:                    nil,
			IPFamilyPolicy:                nil,
			AllocateLoadBalancerNodePorts: nil,
			LoadBalancerClass:             nil,
			InternalTrafficPolicy:         nil,
		},
		Status: v1.ServiceStatus{
			LoadBalancer: v1.LoadBalancerStatus{},
			Conditions:   nil,
		},
	}

	for _, opt := range opts {
		opt(service)
	}

	_, err := c.k8s.CoreV1().Services("default").Create(context.Background(), service, metav1.CreateOptions{})
	require.NoError(t, err)
}
