// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/sync"
	"gorm.io/gorm"
)

const (
	baseTable = "test_grandparents"

	batchAfter = 100

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	log    = logging.LoggerForModule()
	schema = pkgSchema.TestGrandparentsSchema
)

type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)
	Get(ctx context.Context, id string) (*storage.TestGrandparent, bool, error)
	Upsert(ctx context.Context, obj *storage.TestGrandparent) error
	UpsertMany(ctx context.Context, objs []*storage.TestGrandparent) error
	Delete(ctx context.Context, id string) error
	GetIDs(ctx context.Context) ([]string, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.TestGrandparent, []int, error)
	DeleteMany(ctx context.Context, ids []string) error

	Walk(ctx context.Context, fn func(obj *storage.TestGrandparent) error) error

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

func insertIntoTestGrandparents(ctx context.Context, tx pgx.Tx, obj *storage.TestGrandparent) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		obj.GetId(),
		obj.GetVal(),
		serialized,
	}

	finalStr := "INSERT INTO test_grandparents (Id, Val, serialized) VALUES($1, $2, $3) ON CONFLICT(Id) DO UPDATE SET Id = EXCLUDED.Id, Val = EXCLUDED.Val, serialized = EXCLUDED.serialized"
	_, err := tx.Exec(ctx, finalStr, values...)
	if err != nil {
		return err
	}

	var query string

	for childIdx, child := range obj.GetEmbedded() {
		if err := insertIntoTestGrandparentsEmbeddeds(ctx, tx, child, obj.GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from test_grandparents_embeddeds where test_grandparents_Id = $1 AND idx >= $2"
	_, err = tx.Exec(ctx, query, obj.GetId(), len(obj.GetEmbedded()))
	if err != nil {
		return err
	}
	return nil
}

func insertIntoTestGrandparentsEmbeddeds(ctx context.Context, tx pgx.Tx, obj *storage.TestGrandparent_Embedded, test_grandparents_Id string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		test_grandparents_Id,
		idx,
		obj.GetVal(),
	}

	finalStr := "INSERT INTO test_grandparents_embeddeds (test_grandparents_Id, idx, Val) VALUES($1, $2, $3) ON CONFLICT(test_grandparents_Id, idx) DO UPDATE SET test_grandparents_Id = EXCLUDED.test_grandparents_Id, idx = EXCLUDED.idx, Val = EXCLUDED.Val"
	_, err := tx.Exec(ctx, finalStr, values...)
	if err != nil {
		return err
	}

	var query string

	for childIdx, child := range obj.GetEmbedded2() {
		if err := insertIntoTestGrandparentsEmbeddedsEmbedded2(ctx, tx, child, test_grandparents_Id, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from test_grandparents_embeddeds_embedded2 where test_grandparents_Id = $1 AND test_grandparents_embeddeds_idx = $2 AND idx >= $3"
	_, err = tx.Exec(ctx, query, test_grandparents_Id, idx, len(obj.GetEmbedded2()))
	if err != nil {
		return err
	}
	return nil
}

func insertIntoTestGrandparentsEmbeddedsEmbedded2(ctx context.Context, tx pgx.Tx, obj *storage.TestGrandparent_Embedded_Embedded2, test_grandparents_Id string, test_grandparents_embeddeds_idx int, idx int) error {

	values := []interface{}{
		// parent primary keys start
		test_grandparents_Id,
		test_grandparents_embeddeds_idx,
		idx,
		obj.GetVal(),
	}

	finalStr := "INSERT INTO test_grandparents_embeddeds_embedded2 (test_grandparents_Id, test_grandparents_embeddeds_idx, idx, Val) VALUES($1, $2, $3, $4) ON CONFLICT(test_grandparents_Id, test_grandparents_embeddeds_idx, idx) DO UPDATE SET test_grandparents_Id = EXCLUDED.test_grandparents_Id, test_grandparents_embeddeds_idx = EXCLUDED.test_grandparents_embeddeds_idx, idx = EXCLUDED.idx, Val = EXCLUDED.Val"
	_, err := tx.Exec(ctx, finalStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func (s *storeImpl) copyFromTestGrandparents(ctx context.Context, tx pgx.Tx, objs ...*storage.TestGrandparent) error {

	inputRows := [][]interface{}{}

	var err error

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	var deletes []string

	copyCols := []string{

		"id",

		"val",

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

			obj.GetVal(),

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

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"test_grandparents"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err = s.copyFromTestGrandparentsEmbeddeds(ctx, tx, obj.GetId(), obj.GetEmbedded()...); err != nil {
			return err
		}
	}

	return err
}

func (s *storeImpl) copyFromTestGrandparentsEmbeddeds(ctx context.Context, tx pgx.Tx, test_grandparents_Id string, objs ...*storage.TestGrandparent_Embedded) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"test_grandparents_id",

		"idx",

		"val",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{

			test_grandparents_Id,

			idx,

			obj.GetVal(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"test_grandparents_embeddeds"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err = s.copyFromTestGrandparentsEmbeddedsEmbedded2(ctx, tx, test_grandparents_Id, idx, obj.GetEmbedded2()...); err != nil {
			return err
		}
	}

	return err
}

func (s *storeImpl) copyFromTestGrandparentsEmbeddedsEmbedded2(ctx context.Context, tx pgx.Tx, test_grandparents_Id string, test_grandparents_embeddeds_idx int, objs ...*storage.TestGrandparent_Embedded_Embedded2) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"test_grandparents_id",

		"test_grandparents_embeddeds_idx",

		"idx",

		"val",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{

			test_grandparents_Id,

			test_grandparents_embeddeds_idx,

			idx,

			obj.GetVal(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"test_grandparents_embeddeds_embedded2"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

func (s *storeImpl) copyFrom(ctx context.Context, objs ...*storage.TestGrandparent) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "TestGrandparent")
	if err != nil {
		return err
	}
	defer release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	if err := s.copyFromTestGrandparents(ctx, tx, objs...); err != nil {
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

func (s *storeImpl) upsert(ctx context.Context, objs ...*storage.TestGrandparent) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "TestGrandparent")
	if err != nil {
		return err
	}
	defer release()

	for _, obj := range objs {
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}

		if err := insertIntoTestGrandparents(ctx, tx, obj); err != nil {
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

func (s *storeImpl) Upsert(ctx context.Context, obj *storage.TestGrandparent) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Upsert, "TestGrandparent")

	return s.upsert(ctx, obj)
}

func (s *storeImpl) UpsertMany(ctx context.Context, objs []*storage.TestGrandparent) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.UpdateMany, "TestGrandparent")

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
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "TestGrandparent")

	var sacQueryFilter *v1.Query

	return postgres.RunCountRequestForSchema(schema, sacQueryFilter, s.db)
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(ctx context.Context, id string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "TestGrandparent")

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	count, err := postgres.RunCountRequestForSchema(schema, q, s.db)
	return count == 1, err
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(ctx context.Context, id string) (*storage.TestGrandparent, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "TestGrandparent")

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	data, err := postgres.RunGetQueryForSchema(ctx, schema, q, s.db)
	if err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	var msg storage.TestGrandparent
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
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "TestGrandparent")

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	return postgres.RunDeleteRequestForSchema(schema, q, s.db)
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs(ctx context.Context) ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "storage.TestGrandparentIDs")
	var sacQueryFilter *v1.Query

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
func (s *storeImpl) GetMany(ctx context.Context, ids []string) ([]*storage.TestGrandparent, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "TestGrandparent")

	if len(ids) == 0 {
		return nil, nil, nil
	}

	var sacQueryFilter *v1.Query

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
	resultsByID := make(map[string]*storage.TestGrandparent)
	for _, data := range rows {
		msg := &storage.TestGrandparent{}
		if err := proto.Unmarshal(data, msg); err != nil {
			return nil, nil, err
		}
		resultsByID[msg.GetId()] = msg
	}
	missingIndices := make([]int, 0, len(ids)-len(resultsByID))
	// It is important that the elems are populated in the same order as the input ids
	// slice, since some calling code relies on that to maintain order.
	elems := make([]*storage.TestGrandparent, 0, len(resultsByID))
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
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.RemoveMany, "TestGrandparent")

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(ids...).ProtoQuery(),
	)

	return postgres.RunDeleteRequestForSchema(schema, q, s.db)
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(ctx context.Context, fn func(obj *storage.TestGrandparent) error) error {
	var sacQueryFilter *v1.Query
	rows, err := postgres.RunGetManyQueryForSchema(ctx, schema, sacQueryFilter, s.db)
	if err != nil {
		return pgutils.ErrNilIfNoRows(err)
	}
	for _, data := range rows {
		var msg storage.TestGrandparent
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

func dropTableTestGrandparents(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS test_grandparents CASCADE")
	dropTableTestGrandparentsEmbeddeds(ctx, db)

}

func dropTableTestGrandparentsEmbeddeds(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS test_grandparents_embeddeds CASCADE")
	dropTableTestGrandparentsEmbeddedsEmbedded2(ctx, db)

}

func dropTableTestGrandparentsEmbeddedsEmbedded2(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS test_grandparents_embeddeds_embedded2 CASCADE")

}

func Destroy(ctx context.Context, db *pgxpool.Pool) {
	dropTableTestGrandparents(ctx, db)
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
