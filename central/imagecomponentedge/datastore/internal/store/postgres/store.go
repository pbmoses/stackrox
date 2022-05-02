// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/central/metrics"
	pkgSchema "github.com/stackrox/rox/central/postgres/schema"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
)

const (
	baseTable  = "image_component_relations"
	countStmt  = "SELECT COUNT(*) FROM image_component_relations"
	existsStmt = "SELECT EXISTS(SELECT 1 FROM image_component_relations WHERE Id = $1 AND ImageId = $2 AND ImageComponentId = $3 AND ImageComponentName = $4 AND ImageComponentVersion = $5 AND ImageComponentOperatingSystem = $6)"

	getStmt     = "SELECT serialized FROM image_component_relations WHERE Id = $1 AND ImageId = $2 AND ImageComponentId = $3 AND ImageComponentName = $4 AND ImageComponentVersion = $5 AND ImageComponentOperatingSystem = $6"
	deleteStmt  = "DELETE FROM image_component_relations WHERE Id = $1 AND ImageId = $2 AND ImageComponentId = $3 AND ImageComponentName = $4 AND ImageComponentVersion = $5 AND ImageComponentOperatingSystem = $6"
	walkStmt    = "SELECT serialized FROM image_component_relations"
	getIDsStmt  = "SELECT Id FROM image_component_relations"
	getManyStmt = "SELECT serialized FROM image_component_relations WHERE Id = ANY($1::text[])"

	deleteManyStmt = "DELETE FROM image_component_relations WHERE Id = ANY($1::text[])"

	batchAfter = 100

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	log    = logging.LoggerForModule()
	schema = pkgSchema.ImageComponentRelationsSchema
)

type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string, imageId string, imageComponentId string, imageComponentName string, imageComponentVersion string, imageComponentOperatingSystem string) (bool, error)
	Get(ctx context.Context, id string, imageId string, imageComponentId string, imageComponentName string, imageComponentVersion string, imageComponentOperatingSystem string) (*storage.ImageComponentEdge, bool, error)
	GetIDs(ctx context.Context) ([]string, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.ImageComponentEdge, []int, error)

	Walk(ctx context.Context, fn func(obj *storage.ImageComponentEdge) error) error

	AckKeysIndexed(ctx context.Context, keys ...string) error
	GetKeysToIndex(ctx context.Context) ([]string, error)
}

type storeImpl struct {
	db *pgxpool.Pool
}

// New returns a new Store instance using the provided sql instance.
func New(ctx context.Context, db *pgxpool.Pool) Store {
	pgutils.CreateTable(ctx, db, pkgSchema.CreateTableImagesStmt)
	pgutils.CreateTable(ctx, db, pkgSchema.CreateTableImageComponentsStmt)
	pgutils.CreateTable(ctx, db, pkgSchema.CreateTableImageComponentRelationsStmt)

	return &storeImpl{
		db: db,
	}
}

// Count returns the number of objects in the store
func (s *storeImpl) Count(ctx context.Context) (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "ImageComponentEdge")

	row := s.db.QueryRow(ctx, countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(ctx context.Context, id string, imageId string, imageComponentId string, imageComponentName string, imageComponentVersion string, imageComponentOperatingSystem string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "ImageComponentEdge")

	row := s.db.QueryRow(ctx, existsStmt, id, imageId, imageComponentId, imageComponentName, imageComponentVersion, imageComponentOperatingSystem)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, pgutils.ErrNilIfNoRows(err)
	}
	return exists, nil
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(ctx context.Context, id string, imageId string, imageComponentId string, imageComponentName string, imageComponentVersion string, imageComponentOperatingSystem string) (*storage.ImageComponentEdge, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "ImageComponentEdge")

	conn, release, err := s.acquireConn(ctx, ops.Get, "ImageComponentEdge")
	if err != nil {
		return nil, false, err
	}
	defer release()

	row := conn.QueryRow(ctx, getStmt, id, imageId, imageComponentId, imageComponentName, imageComponentVersion, imageComponentOperatingSystem)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	var msg storage.ImageComponentEdge
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

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs(ctx context.Context) ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "storage.ImageComponentEdgeIDs")

	rows, err := s.db.Query(ctx, getIDsStmt)
	if err != nil {
		return nil, pgutils.ErrNilIfNoRows(err)
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// GetMany returns the objects specified by the IDs or the index in the missing indices slice
func (s *storeImpl) GetMany(ctx context.Context, ids []string) ([]*storage.ImageComponentEdge, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "ImageComponentEdge")

	conn, release, err := s.acquireConn(ctx, ops.GetMany, "ImageComponentEdge")
	if err != nil {
		return nil, nil, err
	}
	defer release()

	rows, err := conn.Query(ctx, getManyStmt, ids)
	if err != nil {
		if err == pgx.ErrNoRows {
			missingIndices := make([]int, 0, len(ids))
			for i := range ids {
				missingIndices = append(missingIndices, i)
			}
			return nil, missingIndices, nil
		}
		return nil, nil, err
	}
	defer rows.Close()
	resultsByID := make(map[string]*storage.ImageComponentEdge)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, nil, err
		}
		msg := &storage.ImageComponentEdge{}
		if err := proto.Unmarshal(data, msg); err != nil {
			return nil, nil, err
		}
		resultsByID[msg.GetId()] = msg
	}
	missingIndices := make([]int, 0, len(ids)-len(resultsByID))
	// It is important that the elems are populated in the same order as the input ids
	// slice, since some calling code relies on that to maintain order.
	elems := make([]*storage.ImageComponentEdge, 0, len(resultsByID))
	for i, id := range ids {
		if result, ok := resultsByID[id]; !ok {
			missingIndices = append(missingIndices, i)
		} else {
			elems = append(elems, result)
		}
	}
	return elems, missingIndices, nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(ctx context.Context, fn func(obj *storage.ImageComponentEdge) error) error {
	rows, err := s.db.Query(ctx, walkStmt)
	if err != nil {
		return pgutils.ErrNilIfNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		var msg storage.ImageComponentEdge
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

func dropTableImageComponentRelations(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS image_component_relations CASCADE")

}

func Destroy(ctx context.Context, db *pgxpool.Pool) {
	dropTableImageComponentRelations(ctx, db)
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
