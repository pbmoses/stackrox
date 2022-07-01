// Code generated by pg-bindings generator. DO NOT EDIT.
package n12ton13

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_12_to_n_13_postgres_compliance_operator_scans/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_12_to_n_13_postgres_compliance_operator_scans/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: pkgMigrations.CurrentDBVersionSeqNum() + 12,
		VersionAfter:   storage.Version{SeqNum: int32(pkgMigrations.CurrentDBVersionSeqNum()) + 13},
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.PkgRocksDB, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving compliance_operator_scans from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.ComplianceOperatorScansSchema
	log       = loghelper.LogWrapper{}
)

func move(legacyDB *rocksdb.RocksDB, gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithAllAccess(context.Background())
	store := pgStore.New(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)
	var complianceOperatorScans []*storage.ComplianceOperatorScan
	var err error
	legacyStore.Walk(ctx, func(obj *storage.ComplianceOperatorScan) error {
		complianceOperatorScans = append(complianceOperatorScans, obj)
		if len(complianceOperatorScans) == 10*batchSize {
			if err := store.UpsertMany(ctx, complianceOperatorScans); err != nil {
				log.WriteToStderrf("failed to persist compliance_operator_scans to store %v", err)
				return err
			}
			complianceOperatorScans = complianceOperatorScans[:0]
		}
		return nil
	})
	if len(complianceOperatorScans) > 0 {
		if err = store.UpsertMany(ctx, complianceOperatorScans); err != nil {
			log.WriteToStderrf("failed to persist compliance_operator_scans to store %v", err)
			return err
		}
	}
	return nil
}

func init() {
	migrations.MustRegisterMigration(migration)
}
