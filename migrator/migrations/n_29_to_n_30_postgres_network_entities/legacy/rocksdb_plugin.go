// Code generated by rocksdb-bindings generator. DO NOT EDIT.
package legacy
import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/db"
	"github.com/stackrox/rox/pkg/db/mapcache"
	"github.com/stackrox/rox/pkg/rocksdb"
	generic "github.com/stackrox/rox/pkg/rocksdb/crud"
)

var (
	bucket = []byte("networkentity")
)

type Store interface {
	// Get(ctx context.Context, id string) (*storage.NetworkEntity, bool, error)
    UpsertMany(ctx context.Context, objs []*storage.NetworkEntity) error
	Walk(ctx context.Context, fn func(obj *storage.NetworkEntity) error) error
}

type storeImpl struct {
	crud db.Crud
}

func alloc() proto.Message {
	return &storage.NetworkEntity{}
}

func keyFunc(msg proto.Message) []byte {
	return []byte(msg.(*storage.NetworkEntity).GetInfo().GetId())
}

// New returns a new Store instance using the provided rocksdb instance.
func New(db *rocksdb.RocksDB) (Store, error) {
//	globaldb.RegisterBucket(bucket, "NetworkEntity")
	baseCRUD := generic.NewCRUD(db, bucket, keyFunc, alloc, false)
	cacheCRUD, err := mapcache.NewMapCache(baseCRUD, keyFunc)
	if err != nil {
		return nil, err
	}
	return &storeImpl{
		crud: cacheCRUD,
	}, nil
}
/*
// Get returns the object, if it exists from the store
func (b *storeImpl) Get(_ context.Context, id string) (*storage.NetworkEntity, bool, error) {
	msg, exists, err := b.crud.Get(id)
	if err != nil || !exists {
		return nil, false, err
	}
	return msg.(*storage.NetworkEntity), true, nil
}
*/
// UpsertMany batches objects into the DB
func (b *storeImpl) UpsertMany(_ context.Context, objs []*storage.NetworkEntity) error {
	msgs := make([]proto.Message, 0, len(objs))
	for _, o := range objs {
		msgs = append(msgs, o)
    }

	return b.crud.UpsertMany(msgs)
}
// Walk iterates over all of the objects in the store and applies the closure
func (b *storeImpl) Walk(_ context.Context, fn func(obj *storage.NetworkEntity) error) error {
	return b.crud.Walk(func(msg proto.Message) error {
		return fn(msg.(*storage.NetworkEntity))
	})
}
