// Code generated by pg-bindings generator. DO NOT EDIT.
package n2ton3

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_02_to_n_03_postgres_namespaces/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_02_to_n_03_postgres_namespaces/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: pkgMigrations.CurrentDBVersionSeqNumWithoutPostgres() + 2,
		VersionAfter:   storage.Version{SeqNum: int32(pkgMigrations.CurrentDBVersionSeqNumWithoutPostgres()) + 3},
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving namespaces from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.NamespacesSchema
	log       = loghelper.LogWrapper{}
)

func move(gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithAllAccess(context.Background())
	store := pgStore.New(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)
	var namespaces []*storage.NamespaceMetadata
	var err error
	walk(ctx, legacyStore, func(obj *storage.NamespaceMetadata) error {
		namespaces = append(namespaces, obj)
		if len(namespaces) == batchSize {
			if err := store.UpsertMany(ctx, namespaces); err != nil {
				log.WriteToStderrf("failed to persist namespaces to store %v", err)
				return err
			}
			namespaces = namespaces[:0]
		}
		return nil
	})
	if len(namespaces) > 0 {
		if err = store.UpsertMany(ctx, namespaces); err != nil {
			log.WriteToStderrf("failed to persist namespaces to store %v", err)
			return err
		}
	}
	return nil
}

func walk(ctx context.Context, s legacy.Store, fn func(obj *storage.NamespaceMetadata) error) error {
	return s.Walk(ctx, fn)
}

func init() {
	migrations.MustRegisterMigration(migration)
}
