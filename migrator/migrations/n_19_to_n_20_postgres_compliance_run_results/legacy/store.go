package legacy

import (
	"context"

	"github.com/stackrox/rox/generated/storage"
)

// Store is the interface for accessing stored compliance data
type Store interface {
	Walk(ctx context.Context, fn func(obj *storage.ComplianceRunResults) error) error
	UpsertMany(ctx context.Context, objs []*storage.ComplianceRunResults) error
}
