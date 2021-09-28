package updater

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/activecomponent/converter"
	activeComponent "github.com/stackrox/rox/central/activecomponent/datastore"
	"github.com/stackrox/rox/central/activecomponent/updater/aggregator"
	deploymentStore "github.com/stackrox/rox/central/deployment/datastore"
	imageStore "github.com/stackrox/rox/central/image/datastore"
	processIndicatorStore "github.com/stackrox/rox/central/processindicator/datastore"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/scancomponent"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/simplecache"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	log = logging.LoggerForModule()

	updaterCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.Deployment, resources.Image, resources.Indicator)))
)

type updaterImpl struct {
	acStore         activeComponent.DataStore
	deploymentStore deploymentStore.DataStore
	piStore         processIndicatorStore.DataStore
	imageStore      imageStore.DataStore

	aggregator      aggregator.ProcessAggregator // Aggregator for incoming process indicators
	executableCache simplecache.Cache            // Cache for image scan result

	deploymentToUpdates map[string][]*aggregator.ProcessUpdate
}

type imageExecutable struct {
	execToComponent map[string]string
	scanTime        time.Time
	// Todo(cdu) add scannerVersion string
}

// PopulateExecutableCache extracts executables from image scan and stores them in the executable cache.
// Image executables are cleared on successful return.
func (u *updaterImpl) PopulateExecutableCache(ctx context.Context, image *storage.Image) error {
	if !features.ActiveVulnManagement.Enabled() {
		return nil
	}
	imageID := image.GetId()
	scan := image.GetScan()
	if imageID == "" || scan == nil {
		log.Debugf("no valid scan, skip populating executable cache %s: %s", imageID, image.GetName())
		return nil
	}
	scanTime := protoconv.ConvertTimestampToTimeOrNow(scan.GetScanTime())

	// Check if we should update executable cache
	currRecord, ok := u.executableCache.Get(imageID)
	if ok && !currRecord.(*imageExecutable).scanTime.Before(scanTime) {
		log.Debugf("Skip scan at %v, current scan (%v) has been populated for image %s: %s", scanTime, currRecord.(*imageExecutable).scanTime, imageID, image.GetName())
		return nil
	}

	// Create or update executable cache
	execToComponent := u.getExecToComponentMap(scan)
	u.executableCache.Add(image.GetId(), &imageExecutable{execToComponent: execToComponent, scanTime: scanTime})

	log.Debugf("Executable cache updated for image %s with %d paths", image.GetId(), len(execToComponent))
	return nil
}

func (u *updaterImpl) getExecToComponentMap(imageScan *storage.ImageScan) map[string]string {
	execToComponent := make(map[string]string)

	for _, component := range imageScan.GetComponents() {
		// We do not support non-OS level active components at this time.
		if component.GetSource() != storage.SourceType_OS {
			continue
		}
		componentID := scancomponent.ComponentID(component.GetName(), component.GetVersion())
		for _, exec := range component.GetExecutables() {
			execToComponent[exec.GetPath()] = componentID
		}
		// Remove the executables to save some memory. The same image won't be processed again.
		component.Executables = nil
	}
	return execToComponent
}

// Update detects active components with most recent process run.
func (u *updaterImpl) Update() {
	if !features.ActiveVulnManagement.Enabled() {
		return
	}
	if len(u.deploymentToUpdates) == 0 {
		ids, err := u.deploymentStore.GetDeploymentIDs()
		if err != nil {
			log.Errorf("failed to fetch deployment ids: %v", err)
			return
		}
		u.deploymentToUpdates = u.aggregator.GetAndPrune(u.imageScanned, set.NewStringSet(ids...))
	}

	if err := u.updateActiveComponents(u.deploymentToUpdates); err != nil {
		log.Errorf("failed to update active components: %v", err)
	} else {
		u.deploymentToUpdates = nil
	}

	if err := u.pruneExecutableCache(); err != nil {
		log.Errorf("Error pruning active component executable cache: %v", err)
	}
}

func (u *updaterImpl) imageScanned(imageID string) bool {
	_, ok := u.executableCache.Get(imageID)
	return ok
}

func (u *updaterImpl) updateActiveComponents(deploymentToUpdates map[string][]*aggregator.ProcessUpdate) error {
	for deploymentID, updates := range deploymentToUpdates {
		err := u.updateForDeployment(updaterCtx, deploymentID, updates)
		if err != nil {
			return errors.Wrapf(err, "failed to update active components for deployment %s", deploymentID)
		}
	}
	return nil
}

// updateForDeployment detects and updates active components for a deployment
func (u *updaterImpl) updateForDeployment(ctx context.Context, deploymentID string, updates []*aggregator.ProcessUpdate) error {
	idToContainers := make(map[string]set.StringSet)
	containersToRemove := set.NewStringSet()
	for _, update := range updates {
		if update.ToBeRemoved() {
			containersToRemove.Add(update.ContainerName)
			continue
		}

		if update.FromDatabase() {
			containersToRemove.Add(update.ContainerName)
		}

		result, ok := u.executableCache.Get(update.ImageID)
		if !ok {
			utils.Should(errors.New("cannot find image scan"))
			continue
		}
		execToComponent := result.(*imageExecutable).execToComponent
		execPaths, err := u.getActiveExecPath(deploymentID, update)
		if err != nil {
			return errors.Wrapf(err, "failed to get active executables for deployment %s container %s", deploymentID, update.ContainerName)
		}

		for _, execPath := range execPaths.AsSlice() {
			componentID, ok := execToComponent[execPath]
			if !ok {
				continue
			}
			id := converter.ComposeID(deploymentID, componentID)
			if _, ok := idToContainers[id]; !ok {
				idToContainers[id] = set.NewStringSet()
			}
			containerNameSet := idToContainers[id]
			containerNameSet.Add(update.ContainerName)
		}
	}
	return u.createActiveComponentsAndUpdateDb(ctx, deploymentID, idToContainers, containersToRemove)
}

func (u *updaterImpl) createActiveComponentsAndUpdateDb(ctx context.Context, deploymentID string, acToContexts map[string]set.StringSet, contextsToRemove set.StringSet) error {
	var err error
	var existingAcs []*storage.ActiveComponent
	if contextsToRemove.Cardinality() == 0 {
		var ids []string
		for id := range acToContexts {
			ids = append(ids, id)
		}
		existingAcs, err = u.acStore.GetBatch(ctx, ids)
	} else {
		// Need to check all active components in case there are containers to remove
		query := search.NewQueryBuilder().AddExactMatches(search.DeploymentID, deploymentID).ProtoQuery()
		existingAcs, err = u.acStore.SearchRawActiveComponents(ctx, query)
	}
	if err != nil {
		return err
	}

	var acToRemove []string
	var activeComponents []*converter.CompleteActiveComponent
	for _, ac := range existingAcs {
		updateAc, shouldRemove := merge(ac, contextsToRemove, acToContexts[ac.GetId()])
		if updateAc != nil {
			activeComponents = append(activeComponents, updateAc)
		}
		if shouldRemove {
			acToRemove = append(acToRemove, ac.GetId())
		}
		delete(acToContexts, ac.GetId())
	}
	for id, contexts := range acToContexts {
		_, componentID, err := converter.DecomposeID(id)
		if err != nil {
			utils.Should(err)
			continue
		}
		newAc := &converter.CompleteActiveComponent{
			DeploymentID: deploymentID,
			ComponentID:  componentID,
			ActiveComponent: &storage.ActiveComponent{
				Id:             id,
				ActiveContexts: make(map[string]*storage.ActiveComponent_ActiveContext),
			},
		}
		for context := range contexts {
			newAc.ActiveComponent.ActiveContexts[context] = &storage.ActiveComponent_ActiveContext{ContainerName: context}
		}
		activeComponents = append(activeComponents, newAc)
	}
	log.Debugf("Upserting %d active components and deleting %d for deployment %s", len(activeComponents), len(acToRemove), deploymentID)
	if len(activeComponents) > 0 {
		err = u.acStore.UpsertBatch(ctx, activeComponents)
		if err != nil {
			return errors.Wrapf(err, "failed to upsert %d activeComponents, %v ...", len(activeComponents), activeComponents[:2])
		}
	}
	if len(acToRemove) > 0 {
		err = u.acStore.DeleteBatch(ctx, acToRemove...)
	}
	return err
}

// Merge existing active component with new contexts, addend could be nil
func merge(base *storage.ActiveComponent, subtrahend, addend set.StringSet) (*converter.CompleteActiveComponent, bool) {
	// Only remove the containers that won't be added back.
	toRemove := subtrahend.Difference(addend)
	var changed bool
	for context := range base.ActiveContexts {
		if toRemove.Contains(context) {
			delete(base.ActiveContexts, context)
			changed = true
		}
	}
	for context := range addend {
		if _, ok := base.ActiveContexts[context]; !ok {
			base.ActiveContexts[context] = &storage.ActiveComponent_ActiveContext{ContainerName: context}
			changed = true
		}
	}
	if len(base.ActiveContexts) == 0 {
		return nil, true
	}
	if !changed {
		return nil, false
	}

	deploymentID, componentID, err := converter.DecomposeID(base.GetId())
	if err != nil {
		utils.Should(err)
		return nil, true
	}

	return &converter.CompleteActiveComponent{
		DeploymentID:    deploymentID,
		ComponentID:     componentID,
		ActiveComponent: base,
	}, false
}

func (u *updaterImpl) getActiveExecPath(deploymentID string, update *aggregator.ProcessUpdate) (set.StringSet, error) {
	if update.FromCache() {
		return update.NewPaths, nil
	}
	containerName := update.ContainerName
	query := search.NewQueryBuilder().AddExactMatches(search.DeploymentID, deploymentID).AddExactMatches(search.ContainerName, containerName).ProtoQuery()
	pis, err := u.piStore.SearchRawProcessIndicators(updaterCtx, query)
	if err != nil {
		return nil, err
	}
	execSet := set.NewStringSet()
	for _, pi := range pis {
		if update.ImageID != pi.GetImageId() {
			continue
		}
		execSet.Add(pi.GetSignal().GetExecFilePath())
	}
	log.Debugf("Active executables for %s:%s: %v", deploymentID, containerName, execSet)
	return execSet, nil
}

func (u *updaterImpl) pruneExecutableCache() error {
	results, err := u.imageStore.Search(updaterCtx, search.EmptyQuery())
	if err != nil {
		return err
	}
	imageIDs := search.ResultsToIDSet(results)

	for _, entry := range u.executableCache.Keys() {
		imageID := entry.(string)
		if !imageIDs.Contains(imageID) {
			u.executableCache.Remove(imageID)
		}
	}
	return nil
}
