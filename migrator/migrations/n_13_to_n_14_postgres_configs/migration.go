// Code generated by pg-bindings generator. DO NOT EDIT.
package n13ton14

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_13_to_n_14_postgres_configs/legacy"
	"github.com/stackrox/rox/migrator/types"
	ops "github.com/stackrox/rox/pkg/metrics"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	bolt "go.etcd.io/bbolt"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: 100,
		VersionAfter:   storage.Version{SeqNum: 101},
		Run: func(databases *types.Databases) error {
			legacyStore := legacy.New(databases.BoltDB)
			if err := moveConfigs(databases.BoltDB, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving configs from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.ConfigsSchema
	log       = loghelper.LogWrapper{}
)

func moveConfigs(legacyDB *bolt.DB, gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := context.Background()
	store := newStore(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)

	var configs []*storage.Config
	var err error
	configs, err = legacyStore.GetAll(ctx)
	if err != nil {
		log.WriteToStderr("failed to fetch all configs")
		return err
	}
	if len(configs) > 0 {
		if err = store.copyFrom(ctx, configs...); err != nil {
			log.WriteToStderrf("failed to persist configs to store %v", err)
			return err
		}
	}
	return nil
}

type storeImpl struct {
	db *pgxpool.Pool // Postgres DB
}

// newStore returns a new Store instance using the provided sql instance.
func newStore(db *pgxpool.Pool) *storeImpl {
	return &storeImpl{
		db: db,
	}
}

func (s *storeImpl) acquireConn(ctx context.Context, _ ops.Op, _ string) (*pgxpool.Conn, func(), error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	return conn, conn.Release, nil
}

func init() {
	migrations.MustRegisterMigration(migration)
}
