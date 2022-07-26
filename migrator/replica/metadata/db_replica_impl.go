package metadata

import (
	"github.com/stackrox/rox/pkg/migrations"
)

const (
	// ErrNoPrevious - cannot rollback
	ErrNoPrevious = "Downgrade is not supported. No previous database for force rollback."

	// ErrNoPreviousInDevEnv -- Downgrade is not supported in dev
	ErrNoPreviousInDevEnv = `
Downgrade is not supported.
We compare dev builds by their release tags. For example, 3.0.58.x-58-g848e7365da is greater than
3.0.58.x-57-g848e7365da. However if the dev builds are on diverged branches, the sequence could be wrong.
These builds are not comparable.

To address this:
1. if you are testing migration, you can merge or rebase to make sure the builds are not diverged; or
2. if you simply want to switch the image, you can disable upgrade rollback and bypass this check by:
kubectl -n stackrox set env deploy/central ROX_DONT_COMPARE_DEV_BUILDS=true
`

	// ErrForceUpgradeDisabled -- force rollback is disabled
	ErrForceUpgradeDisabled = "Central force rollback is disabled. If you want to force rollback to the database before last upgrade, please enable force rollback to current version in central config. Note: all data updates since last upgrade will be lost."

	// ErrPreviousMismatchWithVersions -- downgrade is not supported as previous version is too many versions back.
	ErrPreviousMismatchWithVersions = "Database downgrade is not supported. We can only rollback to the central version before last upgrade. Last upgrade %s, current version %s"
)

// DBReplica -- holds information related to DB replicas
type DBReplica struct {
	dirName      string
	migVer       *migrations.MigrationVersion
	databaseName string
}

// GetVersion -- returns the version associated with the replica.
func (d *DBReplica) GetVersion() string {
	if d.migVer == nil {
		return ""
	}
	return d.migVer.MainVersion
}

// GetSeqNum -- returns the sequence number associated with the replica.
func (d *DBReplica) GetSeqNum() int {
	if d.migVer == nil {
		return 0
	}
	return d.migVer.SeqNum
}

// GetDirName -- returns the file system location of the replica.  (Only valid pre-Postgres)
func (d *DBReplica) GetDirName() string {
	return d.dirName
}

// GetDatabaseName -- returns the database name of the replica.  (Postgres)
func (d *DBReplica) GetDatabaseName() string {
	return d.databaseName
}

// GetMigVersion -- returns the migration version associated with the replica.
func (d *DBReplica) GetMigVersion() *migrations.MigrationVersion {
	return d.migVer
}

// New returns a new ready-to-use store.
func New(dirName string, migVer *migrations.MigrationVersion) *DBReplica {
	return &DBReplica{dirName: dirName, migVer: migVer, databaseName: ""}
}

// NewPostgres returns a new ready-to-use store.
func NewPostgres(migVer *migrations.MigrationVersion, databaseName string) *DBReplica {
	return &DBReplica{migVer: migVer, databaseName: databaseName}
}
