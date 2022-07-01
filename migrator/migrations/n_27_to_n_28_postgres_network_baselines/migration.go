// Code generated by pg-bindings generator. DO NOT EDIT.
package n27ton28

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_27_to_n_28_postgres_network_baselines/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_27_to_n_28_postgres_network_baselines/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: pkgMigrations.CurrentDBVersionSeqNum() + 27,
		VersionAfter:   storage.Version{SeqNum: int32(pkgMigrations.CurrentDBVersionSeqNum()) + 28},
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.PkgRocksDB, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving network_baselines from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.NetworkBaselinesSchema
	log       = loghelper.LogWrapper{}
)

func move(legacyDB *rocksdb.RocksDB, gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithAllAccess(context.Background())
	store := pgStore.New(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)
	var networkBaselines []*storage.NetworkBaseline
	var err error
	legacyStore.Walk(ctx, func(obj *storage.NetworkBaseline) error {
		networkBaselines = append(networkBaselines, obj)
		if len(networkBaselines) == 10*batchSize {
			if err := store.UpsertMany(ctx, networkBaselines); err != nil {
				log.WriteToStderrf("failed to persist network_baselines to store %v", err)
				return err
			}
			networkBaselines = networkBaselines[:0]
		}
		return nil
	})
	if len(networkBaselines) > 0 {
		if err = store.UpsertMany(ctx, networkBaselines); err != nil {
			log.WriteToStderrf("failed to persist network_baselines to store %v", err)
			return err
		}
	}
	return nil
}

func init() {
	migrations.MustRegisterMigration(migration)
}
