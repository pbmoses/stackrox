package n30ton31

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_30_to_n_31_postgres_network_flows/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_30_to_n_31_postgres_network_flows/postgres"
	"github.com/stackrox/rox/migrator/migrations/n_30_to_n_31_postgres_network_flows/store"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: pkgMigrations.CurrentDBVersionSeqNum() + 30,
		VersionAfter:   storage.Version{SeqNum: int32(pkgMigrations.CurrentDBVersionSeqNum()) + 31},
		Run: func(databases *types.Databases) error {
			legacyStore := legacy.NewClusterStore(databases.PkgRocksDB)
			if err := move(databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
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

func move(gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore store.ClusterStore) error {
	ctx := sac.WithAllAccess(context.Background())
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)

	clusterStore := pgStore.NewClusterStore(postgresDB)
	var networkBaselines []*storage.NetworkBaseline

	var err error
	walk(ctx, legacyStore, func(obj *storage.NetworkBaseline) error {
		networkBaselines = append(networkBaselines, obj)
		if len(networkBaselines) == batchSize {
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

func walk(ctx context.Context, s store.ClusterStore, fn func(obj *storage.NetworkBaseline) error) error {
	return s.Walk(ctx, fn)
}

func init() {
	migrations.MustRegisterMigration(migration)
}
