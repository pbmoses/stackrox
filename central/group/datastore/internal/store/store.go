package store

import (
	"github.com/stackrox/rox/generated/storage"
)

// Store updates and utilizes groups, which are attribute to role mappings.
//go:generate mockgen-wrapper
type Store interface {
	Get(props *storage.GroupProperties) (*storage.Group, error)
	GetFiltered(func(*storage.GroupProperties) bool) ([]*storage.Group, error)
	GetAll() ([]*storage.Group, error)

	Walk(authProviderID string, attributes map[string][]string) ([]*storage.Group, error)

	Add(*storage.Group) error
	Update(*storage.Group) error
	Mutate(remove, update, add []*storage.Group) error
	Remove(props *storage.GroupProperties) error
}
