package store

// Initially copied with go:generate cp ../../../../central/cve/store/store.go .

import (
	"context"

	"github.com/stackrox/rox/generated/storage"
)

// Store provides storage functionality for CVEs.
type Store interface {
	Count(ctx context.Context) (int, error)
	Get(ctx context.Context, id string) (*storage.CVE, bool, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.CVE, []int, error)
	Walk(ctx context.Context, fn func(obj *storage.CVE) error) error

	Exists(ctx context.Context, id string) (bool, error)

	Upsert(ctx context.Context, cves ...*storage.CVE) error
	Delete(ctx context.Context, ids ...string) error
}
