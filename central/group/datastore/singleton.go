package datastore

import (
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/group/datastore/internal/store/bolt"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	ds   DataStore
	once sync.Once
)

func initialize() {
	if features.PostgresDatastore.Enabled() {
		// nothing yet
	} else {
		ds = New(bolt.New(globaldb.GetGlobalDB()))
	}
}

// Singleton returns the singleton providing access to the roles store.
func Singleton() DataStore {
	once.Do(initialize)
	return ds
}
