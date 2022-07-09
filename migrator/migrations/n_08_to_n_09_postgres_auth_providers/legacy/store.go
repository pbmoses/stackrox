package store

import (
	"context"

	"github.com/stackrox/rox/generated/storage"
)

// Store stores and retrieves providers from the KV storage mechanism.
type Store interface {
	GetAll(ctx context.Context) ([]*storage.AuthProvider, error)

	Exists(ctx context.Context, id string) (bool, error)
	Upsert(ctx context.Context, obj *storage.AuthProvider) error
	Delete(ctx context.Context, id string) error
}
