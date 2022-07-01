package dackbox

import (
	"context"

	"github.com/stackrox/rox/generated/storage"
)

func (b *storeImpl) GetIDs() ([]string, error) {
	dackTxn, err := b.dacky.NewReadOnlyTransaction()
	if err != nil {
		return nil, err
	}
	defer dackTxn.Discard()

	var ids []string
	err = dackTxn.BucketKeyForEach(Bucket, true, func(k []byte) error {
		ids = append(ids, string(k))
		return nil
	})
	return ids, err
}

func (s *storeImpl) Walk(ctx context.Context, fn func(obj *storage.CVE) error) error {
	ids, err := s.GetIDs()
	if err != nil {
		return err
	}

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize

		if end > len(ids) {
			end = len(ids)
		}
		objs, _, err := s.GetMany(ctx, ids[i:end])
		if err != nil {
			return err
		}
		for _, obj := range objs {
			if err = fn(obj); err != nil {
				return err
			}
		}
	}
	return nil
}
