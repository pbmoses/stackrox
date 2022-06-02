// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/sync"
	"gorm.io/gorm"
)

const (
	baseTable = "alerts"

	batchAfter = 100

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	log            = logging.LoggerForModule()
	schema         = pkgSchema.AlertsSchema
	targetResource = resources.Alert
)

type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)
	Get(ctx context.Context, id string) (*storage.Alert, bool, error)
	Upsert(ctx context.Context, obj *storage.Alert) error
	UpsertMany(ctx context.Context, objs []*storage.Alert) error
	Delete(ctx context.Context, id string) error
	GetIDs(ctx context.Context) ([]string, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.Alert, []int, error)
	DeleteMany(ctx context.Context, ids []string) error

	Walk(ctx context.Context, fn func(obj *storage.Alert) error) error

	AckKeysIndexed(ctx context.Context, keys ...string) error
	GetKeysToIndex(ctx context.Context) ([]string, error)
}

type storeImpl struct {
	db    *pgxpool.Pool
	mutex sync.Mutex
}

// New returns a new Store instance using the provided sql instance.
func New(db *pgxpool.Pool) Store {
	return &storeImpl{
		db: db,
	}
}

func insertIntoAlerts(ctx context.Context, tx pgx.Tx, obj *storage.Alert) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		obj.GetId(),
		obj.GetPolicy().GetId(),
		obj.GetPolicy().GetName(),
		obj.GetPolicy().GetDescription(),
		obj.GetPolicy().GetDisabled(),
		obj.GetPolicy().GetCategories(),
		obj.GetPolicy().GetLifecycleStages(),
		obj.GetPolicy().GetSeverity(),
		obj.GetPolicy().GetEnforcementActions(),
		pgutils.NilOrTime(obj.GetPolicy().GetLastUpdated()),
		obj.GetPolicy().GetSORTName(),
		obj.GetPolicy().GetSORTLifecycleStage(),
		obj.GetPolicy().GetSORTEnforcement(),
		obj.GetLifecycleStage(),
		obj.GetClusterId(),
		obj.GetClusterName(),
		obj.GetNamespace(),
		obj.GetNamespaceId(),
		obj.GetDeployment().GetId(),
		obj.GetDeployment().GetName(),
		obj.GetDeployment().GetInactive(),
		obj.GetImage().GetId(),
		obj.GetImage().GetName().GetRegistry(),
		obj.GetImage().GetName().GetRemote(),
		obj.GetImage().GetName().GetTag(),
		obj.GetImage().GetName().GetFullName(),
		obj.GetResource().GetResourceType(),
		obj.GetResource().GetName(),
		obj.GetEnforcement().GetAction(),
		pgutils.NilOrTime(obj.GetTime()),
		obj.GetState(),
		obj.GetTags(),
		serialized,
	}

	finalStr := "INSERT INTO alerts (Id, Policy_Id, Policy_Name, Policy_Description, Policy_Disabled, Policy_Categories, Policy_LifecycleStages, Policy_Severity, Policy_EnforcementActions, Policy_LastUpdated, Policy_SORTName, Policy_SORTLifecycleStage, Policy_SORTEnforcement, LifecycleStage, ClusterId, ClusterName, Namespace, NamespaceId, Deployment_Id, Deployment_Name, Deployment_Inactive, Image_Id, Image_Name_Registry, Image_Name_Remote, Image_Name_Tag, Image_Name_FullName, Resource_ResourceType, Resource_Name, Enforcement_Action, Time, State, Tags, serialized) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33) ON CONFLICT(Id) DO UPDATE SET Id = EXCLUDED.Id, Policy_Id = EXCLUDED.Policy_Id, Policy_Name = EXCLUDED.Policy_Name, Policy_Description = EXCLUDED.Policy_Description, Policy_Disabled = EXCLUDED.Policy_Disabled, Policy_Categories = EXCLUDED.Policy_Categories, Policy_LifecycleStages = EXCLUDED.Policy_LifecycleStages, Policy_Severity = EXCLUDED.Policy_Severity, Policy_EnforcementActions = EXCLUDED.Policy_EnforcementActions, Policy_LastUpdated = EXCLUDED.Policy_LastUpdated, Policy_SORTName = EXCLUDED.Policy_SORTName, Policy_SORTLifecycleStage = EXCLUDED.Policy_SORTLifecycleStage, Policy_SORTEnforcement = EXCLUDED.Policy_SORTEnforcement, LifecycleStage = EXCLUDED.LifecycleStage, ClusterId = EXCLUDED.ClusterId, ClusterName = EXCLUDED.ClusterName, Namespace = EXCLUDED.Namespace, NamespaceId = EXCLUDED.NamespaceId, Deployment_Id = EXCLUDED.Deployment_Id, Deployment_Name = EXCLUDED.Deployment_Name, Deployment_Inactive = EXCLUDED.Deployment_Inactive, Image_Id = EXCLUDED.Image_Id, Image_Name_Registry = EXCLUDED.Image_Name_Registry, Image_Name_Remote = EXCLUDED.Image_Name_Remote, Image_Name_Tag = EXCLUDED.Image_Name_Tag, Image_Name_FullName = EXCLUDED.Image_Name_FullName, Resource_ResourceType = EXCLUDED.Resource_ResourceType, Resource_Name = EXCLUDED.Resource_Name, Enforcement_Action = EXCLUDED.Enforcement_Action, Time = EXCLUDED.Time, State = EXCLUDED.State, Tags = EXCLUDED.Tags, serialized = EXCLUDED.serialized"
	_, err := tx.Exec(ctx, finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func (s *storeImpl) copyFromAlerts(ctx context.Context, tx pgx.Tx, objs ...*storage.Alert) error {

	inputRows := [][]interface{}{}

	var err error

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	var deletes []string

	copyCols := []string{

		"id",

		"policy_id",

		"policy_name",

		"policy_description",

		"policy_disabled",

		"policy_categories",

		"policy_lifecyclestages",

		"policy_severity",

		"policy_enforcementactions",

		"policy_lastupdated",

		"policy_sortname",

		"policy_sortlifecyclestage",

		"policy_sortenforcement",

		"lifecyclestage",

		"clusterid",

		"clustername",

		"namespace",

		"namespaceid",

		"deployment_id",

		"deployment_name",

		"deployment_inactive",

		"image_id",

		"image_name_registry",

		"image_name_remote",

		"image_name_tag",

		"image_name_fullname",

		"resource_resourcetype",

		"resource_name",

		"enforcement_action",

		"time",

		"state",

		"tags",

		"serialized",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		serialized, marshalErr := obj.Marshal()
		if marshalErr != nil {
			return marshalErr
		}

		inputRows = append(inputRows, []interface{}{

			obj.GetId(),

			obj.GetPolicy().GetId(),

			obj.GetPolicy().GetName(),

			obj.GetPolicy().GetDescription(),

			obj.GetPolicy().GetDisabled(),

			obj.GetPolicy().GetCategories(),

			obj.GetPolicy().GetLifecycleStages(),

			obj.GetPolicy().GetSeverity(),

			obj.GetPolicy().GetEnforcementActions(),

			pgutils.NilOrTime(obj.GetPolicy().GetLastUpdated()),

			obj.GetPolicy().GetSORTName(),

			obj.GetPolicy().GetSORTLifecycleStage(),

			obj.GetPolicy().GetSORTEnforcement(),

			obj.GetLifecycleStage(),

			obj.GetClusterId(),

			obj.GetClusterName(),

			obj.GetNamespace(),

			obj.GetNamespaceId(),

			obj.GetDeployment().GetId(),

			obj.GetDeployment().GetName(),

			obj.GetDeployment().GetInactive(),

			obj.GetImage().GetId(),

			obj.GetImage().GetName().GetRegistry(),

			obj.GetImage().GetName().GetRemote(),

			obj.GetImage().GetName().GetTag(),

			obj.GetImage().GetName().GetFullName(),

			obj.GetResource().GetResourceType(),

			obj.GetResource().GetName(),

			obj.GetEnforcement().GetAction(),

			pgutils.NilOrTime(obj.GetTime()),

			obj.GetState(),

			obj.GetTags(),

			serialized,
		})

		// Add the id to be deleted.
		deletes = append(deletes, obj.GetId())

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if err := s.DeleteMany(ctx, deletes); err != nil {
				return err
			}
			// clear the inserts and vals for the next batch
			deletes = nil

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"alerts"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

func (s *storeImpl) copyFrom(ctx context.Context, objs ...*storage.Alert) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "Alert")
	if err != nil {
		return err
	}
	defer release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	if err := s.copyFromAlerts(ctx, tx, objs...); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *storeImpl) upsert(ctx context.Context, objs ...*storage.Alert) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "Alert")
	if err != nil {
		return err
	}
	defer release()

	for _, obj := range objs {
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}

		if err := insertIntoAlerts(ctx, tx, obj); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return err
			}
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *storeImpl) Upsert(ctx context.Context, obj *storage.Alert) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Upsert, "Alert")

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_WRITE_ACCESS).Resource(targetResource).
		ClusterID(obj.GetClusterId()).Namespace(obj.GetNamespace())
	if ok, err := scopeChecker.Allowed(ctx); err != nil {
		return err
	} else if !ok {
		return sac.ErrResourceAccessDenied
	}

	return s.upsert(ctx, obj)
}

func (s *storeImpl) UpsertMany(ctx context.Context, objs []*storage.Alert) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.UpdateMany, "Alert")

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_WRITE_ACCESS).Resource(targetResource)
	if ok, err := scopeChecker.Allowed(ctx); err != nil {
		return err
	} else if !ok {
		var deniedIds []string
		for _, obj := range objs {
			subScopeChecker := scopeChecker.ClusterID(obj.GetClusterId()).Namespace(obj.GetNamespace())
			if ok, err := subScopeChecker.Allowed(ctx); err != nil {
				return err
			} else if !ok {
				deniedIds = append(deniedIds, obj.GetId())
			}
		}
		if len(deniedIds) != 0 {
			return errors.Wrapf(sac.ErrResourceAccessDenied, "modifying alerts with IDs [%s] was denied", strings.Join(deniedIds, ", "))
		}
	}

	// Lock since copyFrom requires a delete first before being executed.  If multiple processes are updating
	// same subset of rows, both deletes could occur before the copyFrom resulting in unique constraint
	// violations
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(objs) < batchAfter {
		return s.upsert(ctx, objs...)
	} else {
		return s.copyFrom(ctx, objs...)
	}
}

// Count returns the number of objects in the store
func (s *storeImpl) Count(ctx context.Context) (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "Alert")

	var sacQueryFilter *v1.Query

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_ACCESS).Resource(targetResource)
	scopeTree, err := scopeChecker.EffectiveAccessScope(permissions.View(targetResource))
	if err != nil {
		return 0, err
	}
	sacQueryFilter, err = sac.BuildClusterNamespaceLevelSACQueryFilter(scopeTree)

	if err != nil {
		return 0, err
	}

	return postgres.RunCountRequestForSchema(schema, sacQueryFilter, s.db)
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(ctx context.Context, id string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "Alert")

	var sacQueryFilter *v1.Query
	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_ACCESS).Resource(targetResource)
	scopeTree, err := scopeChecker.EffectiveAccessScope(permissions.View(targetResource))
	if err != nil {
		return false, err
	}
	sacQueryFilter, err = sac.BuildClusterNamespaceLevelSACQueryFilter(scopeTree)
	if err != nil {
		return false, err
	}

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	count, err := postgres.RunCountRequestForSchema(schema, q, s.db)
	return count == 1, err
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(ctx context.Context, id string) (*storage.Alert, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "Alert")

	var sacQueryFilter *v1.Query

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_ACCESS).Resource(targetResource)
	scopeTree, err := scopeChecker.EffectiveAccessScope(permissions.View(targetResource))
	if err != nil {
		return nil, false, err
	}
	sacQueryFilter, err = sac.BuildClusterNamespaceLevelSACQueryFilter(scopeTree)
	if err != nil {
		return nil, false, err
	}

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	data, err := postgres.RunGetQueryForSchema(ctx, schema, q, s.db)
	if err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	var msg storage.Alert
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, false, err
	}
	return &msg, true, nil
}

func (s *storeImpl) acquireConn(ctx context.Context, op ops.Op, typ string) (*pgxpool.Conn, func(), error) {
	defer metrics.SetAcquireDBConnDuration(time.Now(), op, typ)
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	return conn, conn.Release, nil
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(ctx context.Context, id string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "Alert")

	var sacQueryFilter *v1.Query
	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_WRITE_ACCESS).Resource(targetResource)
	scopeTree, err := scopeChecker.EffectiveAccessScope(permissions.Modify(targetResource))
	if err != nil {
		return err
	}
	sacQueryFilter, err = sac.BuildClusterNamespaceLevelSACQueryFilter(scopeTree)
	if err != nil {
		return err
	}

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	return postgres.RunDeleteRequestForSchema(schema, q, s.db)
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs(ctx context.Context) ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "storage.AlertIDs")
	var sacQueryFilter *v1.Query

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_ACCESS).Resource(targetResource)
	scopeTree, err := scopeChecker.EffectiveAccessScope(permissions.View(targetResource))
	if err != nil {
		return nil, err
	}
	sacQueryFilter, err = sac.BuildClusterNamespaceLevelSACQueryFilter(scopeTree)
	if err != nil {
		return nil, err
	}
	result, err := postgres.RunSearchRequestForSchema(schema, sacQueryFilter, s.db)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(result))
	for _, entry := range result {
		ids = append(ids, entry.ID)
	}

	return ids, nil
}

// GetMany returns the objects specified by the IDs or the index in the missing indices slice
func (s *storeImpl) GetMany(ctx context.Context, ids []string) ([]*storage.Alert, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "Alert")

	if len(ids) == 0 {
		return nil, nil, nil
	}

	var sacQueryFilter *v1.Query

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_ACCESS).Resource(targetResource)
	scopeTree, err := scopeChecker.EffectiveAccessScope(permissions.ResourceWithAccess{
		Resource: targetResource,
		Access:   storage.Access_READ_ACCESS,
	})
	if err != nil {
		return nil, nil, err
	}
	sacQueryFilter, err = sac.BuildClusterNamespaceLevelSACQueryFilter(scopeTree)
	if err != nil {
		return nil, nil, err
	}
	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(ids...).ProtoQuery(),
	)

	rows, err := postgres.RunGetManyQueryForSchema(ctx, schema, q, s.db)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			missingIndices := make([]int, 0, len(ids))
			for i := range ids {
				missingIndices = append(missingIndices, i)
			}
			return nil, missingIndices, nil
		}
		return nil, nil, err
	}
	resultsByID := make(map[string]*storage.Alert)
	for _, data := range rows {
		msg := &storage.Alert{}
		if err := proto.Unmarshal(data, msg); err != nil {
			return nil, nil, err
		}
		resultsByID[msg.GetId()] = msg
	}
	missingIndices := make([]int, 0, len(ids)-len(resultsByID))
	// It is important that the elems are populated in the same order as the input ids
	// slice, since some calling code relies on that to maintain order.
	elems := make([]*storage.Alert, 0, len(resultsByID))
	for i, id := range ids {
		if result, ok := resultsByID[id]; !ok {
			missingIndices = append(missingIndices, i)
		} else {
			elems = append(elems, result)
		}
	}
	return elems, missingIndices, nil
}

// Delete removes the specified IDs from the store
func (s *storeImpl) DeleteMany(ctx context.Context, ids []string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.RemoveMany, "Alert")

	var sacQueryFilter *v1.Query

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_WRITE_ACCESS).Resource(targetResource)
	scopeTree, err := scopeChecker.EffectiveAccessScope(permissions.Modify(targetResource))
	if err != nil {
		return err
	}
	sacQueryFilter, err = sac.BuildClusterNamespaceLevelSACQueryFilter(scopeTree)
	if err != nil {
		return err
	}

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(ids...).ProtoQuery(),
	)

	return postgres.RunDeleteRequestForSchema(schema, q, s.db)
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(ctx context.Context, fn func(obj *storage.Alert) error) error {
	var sacQueryFilter *v1.Query
	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_ACCESS).Resource(targetResource)
	scopeTree, err := scopeChecker.EffectiveAccessScope(permissions.ResourceWithAccess{
		Resource: targetResource,
		Access:   storage.Access_READ_ACCESS,
	})
	if err != nil {
		return err
	}
	sacQueryFilter, err = sac.BuildClusterNamespaceLevelSACQueryFilter(scopeTree)
	if err != nil {
		return err
	}
	rows, err := postgres.RunGetManyQueryForSchema(ctx, schema, sacQueryFilter, s.db)
	if err != nil {
		return pgutils.ErrNilIfNoRows(err)
	}
	for _, data := range rows {
		var msg storage.Alert
		if err := proto.Unmarshal(data, &msg); err != nil {
			return err
		}
		if err := fn(&msg); err != nil {
			return err
		}
	}
	return nil
}

//// Used for testing

func dropTableAlerts(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS alerts CASCADE")

}

func Destroy(ctx context.Context, db *pgxpool.Pool) {
	dropTableAlerts(ctx, db)
}

// NewTestStore returns a new Store instance for testing
func NewTestStore(ctx context.Context, db *pgxpool.Pool, gormDB *gorm.DB) Store {
	pkgSchema.ApplySchemaForTable(ctx, gormDB, baseTable)
	return New(db)
}

//// Stubs for satisfying legacy interfaces

// AckKeysIndexed acknowledges the passed keys were indexed
func (s *storeImpl) AckKeysIndexed(ctx context.Context, keys ...string) error {
	return nil
}

// GetKeysToIndex returns the keys that need to be indexed
func (s *storeImpl) GetKeysToIndex(ctx context.Context) ([]string, error) {
	return nil, nil
}
