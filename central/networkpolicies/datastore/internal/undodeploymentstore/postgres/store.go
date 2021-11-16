// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/batcher"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/set"
)

const (
		countStmt = "select count(*) from networkpolicies-undodeployment"
		existsStmt = "select exists(select 1 from networkpolicies-undodeployment where id = $1)"
		getIDsStmt = "select id from networkpolicies-undodeployment"
		getStmt = "select value from networkpolicies-undodeployment where id = $1"
		getManyStmt = "select value from networkpolicies-undodeployment where id = ANY($1::text[])"
		upsertStmt = "insert into networkpolicies-undodeployment (id, value) values($1, $2) on conflict(id) do update set value = EXCLUDED.value"
		deleteStmt = "delete from networkpolicies-undodeployment where id = $1"
		deleteManyStmt = "delete from networkpolicies-undodeployment where id = ANY($1::text[])"
		walkStmt = "select value from networkpolicies-undodeployment"
		walkWithIDStmt = "select id, value from networkpolicies-undodeployment"
)

var (
	log = logging.LoggerForModule()

	table = "networkpolicies-undodeployment"

	marshaler = &jsonpb.Marshaler{EnumsAsInts: true, EmitDefaults: true}
)

type Store interface {
	Count() (int, error)
	Exists(id string) (bool, error)
	GetIDs() ([]string, error)
	Get(id string) (*storage.NetworkPolicyApplicationUndoDeploymentRecord, bool, error)
	GetMany(ids []string) ([]*storage.NetworkPolicyApplicationUndoDeploymentRecord, []int, error)
	Upsert(obj *storage.NetworkPolicyApplicationUndoDeploymentRecord) error
	UpsertMany(objs []*storage.NetworkPolicyApplicationUndoDeploymentRecord) error
	Delete(id string) error
	DeleteMany(ids []string) error
	Walk(fn func(obj *storage.NetworkPolicyApplicationUndoDeploymentRecord) error) error
	AckKeysIndexed(keys ...string) error
	GetKeysToIndex() ([]string, error)
}

type storeImpl struct {
	db *pgxpool.Pool
}

func alloc() proto.Message {
	return &storage.NetworkPolicyApplicationUndoDeploymentRecord{}
}

func keyFunc(msg proto.Message) string {
	return msg.(*storage.NetworkPolicyApplicationUndoDeploymentRecord).GetDeploymentId()
}

const (
	createTableQuery = "create table if not exists networkpolicies-undodeployment (id varchar primary key, value jsonb)"
	createIDIndexQuery = "create index if not exists networkpolicies-undodeployment_id on networkpolicies-undodeployment using hash ((id))"

	batchInsertTemplate = "insert into networkpolicies-undodeployment (id, value) values %s on conflict(id) do update set value = EXCLUDED.value"
)

// New returns a new Store instance using the provided sql instance.
func New(db *pgxpool.Pool) Store {
	globaldb.RegisterTable(table, "NetworkPolicyApplicationUndoDeploymentRecord")

	for _, table := range []string {
		"create table if not exists NetworkPolicyApplicationUndoDeploymentRecord(serialized jsonb not null, PRIMARY KEY ());",
		"create table if not exists NetworkPolicyApplicationUndoDeploymentRecord_ToDelete(idx numeric not null, PRIMARY KEY (idx), CONSTRAINT fk_parent_table FOREIGN KEY () REFERENCES NetworkPolicyApplicationUndoDeploymentRecord() ON DELETE CASCADE);",
		"create table if not exists NetworkPolicyApplicationUndoDeploymentRecord_ToDelete(idx numeric not null, PRIMARY KEY (idx), CONSTRAINT fk_parent_table FOREIGN KEY () REFERENCES NetworkPolicyApplicationUndoDeploymentRecord() ON DELETE CASCADE);",
		
	} {
		_, err := db.Exec(context.Background(), table)
		if err != nil {
			panic("error creating table: " + table)
		}
	}

	// Will also autogen the indexes in the future
	//_, err := db.Exec(context.Background(), createIDIndexQuery)
	//if err != nil {
	//	panic("error creating index")
	//}

//
	return &storeImpl{
		db: db,
	}
//
}

// Count returns the number of objects in the store
func (s *storeImpl) Count() (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "NetworkPolicyApplicationUndoDeploymentRecord")

	row := s.db.QueryRow(context.Background(), countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(id string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "NetworkPolicyApplicationUndoDeploymentRecord")

	row := s.db.QueryRow(context.Background(), existsStmt, id)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, nilNoRows(err)
	}
	return exists, nil
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs() ([]string, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetAll, "NetworkPolicyApplicationUndoDeploymentRecordIDs")

	rows, err := s.db.Query(context.Background(), getIDsStmt)
	if err != nil {
		return nil, nilNoRows(err)
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

func nilNoRows(err error) error {
	if err == pgx.ErrNoRows {
		return nil
	}
	return err
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(id string) (*storage.NetworkPolicyApplicationUndoDeploymentRecord, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "NetworkPolicyApplicationUndoDeploymentRecord")

	conn, release := s.acquireConn(ops.Get, "NetworkPolicyApplicationUndoDeploymentRecord")
	defer release()

	row := conn.QueryRow(context.Background(), getStmt, id)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, nilNoRows(err)
	}

	msg := alloc()
	buf := bytes.NewBuffer(data)
	defer metrics.SetJSONPBOperationDurationTime(time.Now(), "Unmarshal", "NetworkPolicyApplicationUndoDeploymentRecord")
	if err := jsonpb.Unmarshal(buf, msg); err != nil {
		return nil, false, err
	}
	return msg.(*storage.NetworkPolicyApplicationUndoDeploymentRecord), true, nil
}

// GetMany returns the objects specified by the IDs or the index in the missing indices slice 
func (s *storeImpl) GetMany(ids []string) ([]*storage.NetworkPolicyApplicationUndoDeploymentRecord, []int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.GetMany, "NetworkPolicyApplicationUndoDeploymentRecord")

	conn, release := s.acquireConn(ops.GetMany, "NetworkPolicyApplicationUndoDeploymentRecord")
	defer release()

	rows, err := conn.Query(context.Background(), getManyStmt, ids)
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
	elems := make([]*storage.NetworkPolicyApplicationUndoDeploymentRecord, 0, len(ids))
	foundSet := set.NewStringSet()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, nil, err
		}
		msg := alloc()
		buf := bytes.NewBuffer(data)
		t := time.Now()
		if err := jsonpb.Unmarshal(buf, msg); err != nil {
			return nil, nil, err
		}
		metrics.SetJSONPBOperationDurationTime(t, "Unmarshal", "NetworkPolicyApplicationUndoDeploymentRecord")
		elem := msg.(*storage.NetworkPolicyApplicationUndoDeploymentRecord)
		foundSet.Add(elem.GetId())
		elems = append(elems, elem)
	}
	missingIndices := make([]int, 0, len(ids)-len(foundSet))
	for i, id := range ids {
		if !foundSet.Contains(id) {
			missingIndices = append(missingIndices, i)
		}
	}
	return elems, missingIndices, nil
}

func nilOrStringTimestamp(t *types.Timestamp) *string {
  if t == nil {
    return nil
  }
  s := t.String()
  return &s
}

func (s *storeImpl) upsert(id string, obj0 *storage.NetworkPolicyApplicationUndoDeploymentRecord) error {
	t := time.Now()
	serialized, err := marshaler.MarshalToString(obj0)
	if err != nil {
		return err
	}
	metrics.SetJSONPBOperationDurationTime(t, "Marshal", "NetworkPolicyApplicationUndoDeploymentRecord")
	conn, release := s.acquireConn(ops.Add, "NetworkPolicyApplicationUndoDeploymentRecord")
	defer release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit(context.Background())
		}
//else {
//			if rollBackErr := tx.Rollback(context.Background()); rollBackErr != nil {
//				// multi error?
//				err = rollBackErr
//			}
//		}
	}()

	localQuery := "insert into NetworkPolicyApplicationUndoDeploymentRecord(serialized) values($1) on conflict() do update set serialized = EXCLUDED.serialized"
_, err = tx.Exec(context.Background(), localQuery, serialized)
if err != nil {
    return err
  }


	return err
}

// Upsert inserts the object into the DB
func (s *storeImpl) Upsert(obj *storage.NetworkPolicyApplicationUndoDeploymentRecord) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Add, "NetworkPolicyApplicationUndoDeploymentRecord")
	return s.upsert(keyFunc(obj), obj)
}

func (s *storeImpl) acquireConn(op ops.Op, typ string) (*pgxpool.Conn, func()) {
	defer metrics.SetAcquireDuration(time.Now(), op, typ)
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		panic(err)
	}
	return conn, conn.Release
}

// UpsertMany batches objects into the DB
func (s *storeImpl) UpsertMany(objs []*storage.NetworkPolicyApplicationUndoDeploymentRecord) error {
	if len(objs) == 0 {
		return nil
	}

	conn, release := s.acquireConn(ops.AddMany, "NetworkPolicyApplicationUndoDeploymentRecord")
	defer release()

	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.AddMany, "NetworkPolicyApplicationUndoDeploymentRecord")
	numElems := 2
	batch := batcher.New(len(objs), 60000/numElems)
	for start, end, ok := batch.Next(); ok; start, end, ok = batch.Next() {
		var placeholderStr string
		data := make([]interface{}, 0, numElems * len(objs))
		for i, obj := range objs[start:end] {
			if i != 0 {
				placeholderStr += ", "
			}
			placeholderStr += postgres.GetValues(i*numElems+1, (i+1)*numElems+1)

			t := time.Now()
			value, err := marshaler.MarshalToString(obj)
			if err != nil {
				return err
			}
			metrics.SetJSONPBOperationDurationTime(t, "Marshal", "NetworkPolicyApplicationUndoDeploymentRecord")
			id := keyFunc(obj)
			data = append(data, id, value)
		}
		if _, err := conn.Exec(context.Background(), fmt.Sprintf(batchInsertTemplate, placeholderStr), data...); err != nil {
			return err
		}
	}
	return nil
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(id string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Remove, "NetworkPolicyApplicationUndoDeploymentRecord")

	conn, release := s.acquireConn(ops.Remove, "NetworkPolicyApplicationUndoDeploymentRecord")
	defer release()

	if _, err := conn.Exec(context.Background(), deleteStmt, id); err != nil {
		return err
	}
	return nil
}

// Delete removes the specified IDs from the store
func (s *storeImpl) DeleteMany(ids []string) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.RemoveMany, "NetworkPolicyApplicationUndoDeploymentRecord")

	conn, release := s.acquireConn(ops.RemoveMany, "NetworkPolicyApplicationUndoDeploymentRecord")
	defer release()
	if _, err := conn.Exec(context.Background(), deleteManyStmt, ids); err != nil {
		return err
	}
	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(fn func(obj *storage.NetworkPolicyApplicationUndoDeploymentRecord) error) error {
	rows, err := s.db.Query(context.Background(), walkStmt)
	if err != nil {
		return nilNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		msg := alloc()
		buf := bytes.NewReader(data)
		if err := jsonpb.Unmarshal(buf, msg); err != nil {
			return err
		}
		return fn(msg.(*storage.NetworkPolicyApplicationUndoDeploymentRecord))
	}
	return nil
}

// AckKeysIndexed acknowledges the passed keys were indexed
func (s *storeImpl) AckKeysIndexed(keys ...string) error {
	return nil
}

// GetKeysToIndex returns the keys that need to be indexed
func (s *storeImpl) GetKeysToIndex() ([]string, error) {
	return nil, nil
}
