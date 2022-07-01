// Code generated by pg-bindings generator. DO NOT EDIT.
package n1ton2

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_01_to_n_02_postgres_alerts/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_01_to_n_02_postgres_alerts/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: pkgMigrations.CurrentDBVersionSeqNum() + 1,
		VersionAfter:   storage.Version{SeqNum: int32(pkgMigrations.CurrentDBVersionSeqNum()) + 2},
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.PkgRocksDB, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving alerts from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.AlertsSchema
	log       = loghelper.LogWrapper{}
)

func move(legacyDB *rocksdb.RocksDB, gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowAllAccessScopeChecker())
	store := pgStore.New(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)
	var alerts []*storage.Alert
	var err error
	legacyStore.Walk(ctx, func(obj *storage.Alert) error {
		alerts = append(alerts, obj)
		if len(alerts) == 10*batchSize {
			if err := store.UpsertMany(ctx, alerts); err != nil {
				log.WriteToStderrf("failed to persist alerts to store %v", err)
				return err
			}
			alerts = alerts[:0]
		}
		return nil
	})
	if len(alerts) > 0 {
		if err = store.UpsertMany(ctx, alerts); err != nil {
			log.WriteToStderrf("failed to persist alerts to store %v", err)
			return err
		}
	}
	return nil
}

func init() {
	migrations.MustRegisterMigration(migration)
}
