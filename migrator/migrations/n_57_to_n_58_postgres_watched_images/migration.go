// Code generated by pg-bindings generator. DO NOT EDIT.
package n57ton58

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_57_to_n_58_postgres_watched_images/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_57_to_n_58_postgres_watched_images/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: pkgMigrations.CurrentDBVersionSeqNum() + 57,
		VersionAfter:   storage.Version{SeqNum: int32(pkgMigrations.CurrentDBVersionSeqNum()) + 58},
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.PkgRocksDB, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving watched_images from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.WatchedImagesSchema
	log       = loghelper.LogWrapper{}
)

func move(legacyDB *rocksdb.RocksDB, gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithAllAccess(context.Background())
	store := pgStore.New(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)
	var watchedImages []*storage.WatchedImage
	var err error
	legacyStore.Walk(ctx, func(obj *storage.WatchedImage) error {
		watchedImages = append(watchedImages, obj)
		if len(watchedImages) == 10*batchSize {
			if err := store.UpsertMany(ctx, watchedImages); err != nil {
				log.WriteToStderrf("failed to persist watched_images to store %v", err)
				return err
			}
			watchedImages = watchedImages[:0]
		}
		return nil
	})
	if len(watchedImages) > 0 {
		if err = store.UpsertMany(ctx, watchedImages); err != nil {
			log.WriteToStderrf("failed to persist watched_images to store %v", err)
			return err
		}
	}
	return nil
}

func init() {
	migrations.MustRegisterMigration(migration)
}
