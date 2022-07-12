// Code generated by pg-bindings generator. DO NOT EDIT.
package n32ton33

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_32_to_n_33_postgres_networkpolicies/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_32_to_n_33_postgres_networkpolicies/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: pkgMigrations.CurrentDBVersionSeqNumWithoutPostgres() + 32,
		VersionAfter:   storage.Version{SeqNum: int32(pkgMigrations.CurrentDBVersionSeqNumWithoutPostgres()) + 33},
		Run: func(databases *types.Databases) error {
			legacyStore := legacy.New(databases.BoltDB)
			if err := move(databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving networkpolicies from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.NetworkpoliciesSchema
	log       = loghelper.LogWrapper{}
)

func move(gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithAllAccess(context.Background())
	store := pgStore.New(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)
	var networkpolicies []*storage.NetworkPolicy
	var err error
	walk(ctx, legacyStore, func(obj *storage.NetworkPolicy) error {
		networkpolicies = append(networkpolicies, obj)
		if len(networkpolicies) == batchSize {
			if err := store.UpsertMany(ctx, networkpolicies); err != nil {
				log.WriteToStderrf("failed to persist networkpolicies to store %v", err)
				return err
			}
			networkpolicies = networkpolicies[:0]
		}
		return nil
	})
	if len(networkpolicies) > 0 {
		if err = store.UpsertMany(ctx, networkpolicies); err != nil {
			log.WriteToStderrf("failed to persist networkpolicies to store %v", err)
			return err
		}
	}
	return nil
}

func walk(ctx context.Context, s legacy.Store, fn func(obj *storage.NetworkPolicy) error) error {
	return s.Walk(ctx, fn)
}

func init() {
	migrations.MustRegisterMigration(migration)
}
