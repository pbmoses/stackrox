package n39ton40

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_39_to_n_40_postgres_policies/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_39_to_n_40_postgres_policies/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	bolt "go.etcd.io/bbolt"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: 100,
		VersionAfter:   storage.Version{SeqNum: 101},
		Run: func(databases *types.Databases) error {
			legacyStore := legacy.New(databases.BoltDB)
			if err := move(databases.BoltDB, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving policies from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.PoliciesSchema
	log       = loghelper.LogWrapper{}
)

func move(legacyDB *bolt.DB, gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowAllAccessScopeChecker())
	store := pgStore.New(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)
	var policies []*storage.Policy
	var err error
	legacyStore.Walk(ctx, func(obj *storage.Policy) error {
		policies = append(policies, obj)
		if len(policies) == 10*batchSize {
			if err := store.UpsertMany(ctx, policies); err != nil {
				log.WriteToStderrf("failed to persist policies to store %v", err)
				return err
			}
			policies = policies[:0]
		}
		return nil
	})
	if len(policies) > 0 {
		if err = store.UpsertMany(ctx, policies); err != nil {
			log.WriteToStderrf("failed to persist policies to store %v", err)
			return err
		}
	}
	return nil
}

func init() {
	migrations.MustRegisterMigration(migration)
}
