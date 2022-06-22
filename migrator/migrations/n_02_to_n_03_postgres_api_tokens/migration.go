// Code generated by pg-bindings generator. DO NOT EDIT.
package n2ton3

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_02_to_n_03_postgres_api_tokens/legacy"
	"github.com/stackrox/rox/migrator/types"
	ops "github.com/stackrox/rox/pkg/metrics"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: 100,
		VersionAfter:   storage.Version{SeqNum: 101},
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := moveApiTokens(databases.PkgRocksDB, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving api_tokens from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.ApiTokensSchema
	log       = loghelper.LogWrapper{}
)

func moveApiTokens(legacyDB *rocksdb.RocksDB, gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := context.Background()
	store := newStore(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)

	var apiTokens []*storage.TokenMetadata
	var err error
	legacyStore.Walk(ctx, func(obj *storage.TokenMetadata) error {
		apiTokens = append(apiTokens, obj)
		if len(apiTokens) == 10*batchSize {
			if err := store.copyFrom(ctx, apiTokens...); err != nil {
				log.WriteToStderrf("failed to persist api_tokens to store %v", err)
				return err
			}
			apiTokens = apiTokens[:0]
		}
		return nil
	})
	if len(apiTokens) > 0 {
		if err = store.copyFrom(ctx, apiTokens...); err != nil {
			log.WriteToStderrf("failed to persist api_tokens to store %v", err)
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
func (s *storeImpl) DeleteMany(ctx context.Context, ids []string) error {
	q := search.NewQueryBuilder().AddDocIDs(ids...).ProtoQuery()
	return postgres.RunDeleteRequestForSchema(schema, q, s.db)
}

func init() {
	migrations.MustRegisterMigration(migration)
}
