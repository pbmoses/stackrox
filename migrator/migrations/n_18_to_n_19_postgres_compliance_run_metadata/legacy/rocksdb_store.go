package legacy

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/dbhelper"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/rocksdb"
	generic "github.com/stackrox/rox/pkg/rocksdb/crud"
	"github.com/stackrox/rox/pkg/timestamp"
	"github.com/tecbot/gorocksdb"
)

var (
	readOptions  = generic.DefaultReadOptions()
	writeOptions = generic.DefaultWriteOptions()

	resultsBucketName = []byte("compliance-run-results")

	metadataKey = dbhelper.GetBucketKey(resultsBucketName, []byte("metadata"))

	log = logging.LoggerForModule()
)

// New returns a compliance results store that is backed by RocksDB.
func New(db *rocksdb.RocksDB) (Store, error) {
	// globaldb.RegisterBucket(resultsBucketName, "ComplianceRunResults")
	return &rocksdbStore{
		db: db,
	}, nil
}

type rocksdbStore struct {
	db *rocksdb.RocksDB
}

func (r *rocksdbStore) Walk(ctx context.Context, fn func(obj *storage.ComplianceRunMetadata) error) error {
	iterator := r.db.NewIterator(readOptions)
	defer iterator.Close()
	// Runs are sorted by time so we must iterate over each key to see if it has the correct run ID.
	for iterator.Seek(metadataKey); iterator.ValidForPrefix(metadataKey); iterator.Next() {
		domain, err := unmarshalMetadata(iterator)
		if err != nil {
			return err
		}
		if err = fn(domain); err != nil {
			return err
		}
	}
	return nil
}

func (r *rocksdbStore) UpsertMany(ctx context.Context, objs []*storage.ComplianceRunMetadata) error {
	for _, obj := range objs {
		if err := r.StoreComplianceRunMetadata(ctx, obj); err != nil {
			return err
		}
	}
	return nil
}

func unmarshalMetadata(iterator *gorocksdb.Iterator) (*storage.ComplianceRunMetadata, error) {
	bytes := iterator.Value().Data()
	if len(bytes) == 0 {
		return nil, errors.New("compliance domain data is empty")
	}
	var domain storage.ComplianceRunMetadata
	if err := domain.Unmarshal(bytes); err != nil {
		return nil, errors.Wrap(err, "unmarshalling compliance domain")
	}
	return &domain, nil
}

type keyMaker struct {
	partialMetadataPrefix []byte
}

func (k *keyMaker) getMetadataIterationPrefix() []byte {
	return k.partialMetadataPrefix
}

func (k *keyMaker) getKeysForMetadata(metadata *storage.ComplianceRunMetadata) ([]byte, error) {
	runID := metadata.GetRunId()
	if runID == "" {
		return nil, errors.New("run has an empty ID")
	}

	tsBytes := []byte(fmt.Sprintf("%016X", timestamp.FromGoTime(time.Now())))
	// Invert the bits of each byte of the timestamp in order to have the most recent timestamp first
	for i, tsByte := range tsBytes {
		tsBytes[i] = -tsByte
	}
	separatorAndRunID := []byte(fmt.Sprintf(":%s", runID))
	tsAndRunIDPrefix := append(tsBytes, separatorAndRunID...)

	key := append([]byte{}, k.partialMetadataPrefix...)
	key = append(metadataKey, tsAndRunIDPrefix...)

	return key, nil
}

func getPrefix(leftPrefix, rightPrefix string) []byte {
	return []byte(leftPrefix + ":" + rightPrefix)
}

func getClusterStandardPrefixes(clusterID, standardID string) []byte {
	// trailing colon is intentional, this prefix will always be followed by a timestamp and a run ID
	partialPrefix := fmt.Sprintf("%s:%s:", clusterID, standardID)
	metadataPrefix := getPrefix(string(metadataKey), partialPrefix)
	return metadataPrefix
}

func getKeyMaker(clusterID, standardID string) *keyMaker {
	metadataPrefix := getClusterStandardPrefixes(clusterID, standardID)

	return &keyMaker{
		partialMetadataPrefix: metadataPrefix,
	}
}

func (r *rocksdbStore) StoreComplianceRunMetadata(_ context.Context, metadata *storage.ComplianceRunMetadata) error {
	clusterID := metadata.GetClusterId()
	standardID := metadata.GetStandardId()

	serializedMD, err := metadata.Marshal()
	if err != nil {
		return errors.Wrap(err, "serializing metadata")
	}

	keyMaker := getKeyMaker(clusterID, standardID)
	key, err := keyMaker.getKeysForMetadata(metadata)
	if err != nil {
		return err
	}

	err = r.db.Put(writeOptions, key, serializedMD)
	return errors.Wrap(err, "storing metadata")
}
