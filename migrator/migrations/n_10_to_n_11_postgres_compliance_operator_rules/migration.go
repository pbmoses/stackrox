// Code generated by pg-bindings generator. DO NOT EDIT.
package n10ton11

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_10_to_n_11_postgres_compliance_operator_rules/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_10_to_n_11_postgres_compliance_operator_rules/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

var (
	migration = types.Migration{
		StartingSeqNum: pkgMigrations.CurrentDBVersionSeqNum() + 10,
		VersionAfter:   storage.Version{SeqNum: int32(pkgMigrations.CurrentDBVersionSeqNum()) + 11},
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.PkgRocksDB, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving compliance_operator_rules from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = pkgSchema.ComplianceOperatorRulesSchema
	log       = loghelper.LogWrapper{}
)

func move(legacyDB *rocksdb.RocksDB, gormDB *gorm.DB, postgresDB *pgxpool.Pool, legacyStore legacy.Store) error {
	ctx := sac.WithAllAccess(context.Background())
	store := pgStore.New(postgresDB)
	pkgSchema.ApplySchemaForTable(context.Background(), gormDB, schema.Table)
	var complianceOperatorRules []*storage.ComplianceOperatorRule
	var err error
	legacyStore.Walk(ctx, func(obj *storage.ComplianceOperatorRule) error {
		complianceOperatorRules = append(complianceOperatorRules, obj)
		if len(complianceOperatorRules) == 10*batchSize {
			if err := store.UpsertMany(ctx, complianceOperatorRules); err != nil {
				log.WriteToStderrf("failed to persist compliance_operator_rules to store %v", err)
				return err
			}
			complianceOperatorRules = complianceOperatorRules[:0]
		}
		return nil
	})
	if len(complianceOperatorRules) > 0 {
		if err = store.UpsertMany(ctx, complianceOperatorRules); err != nil {
			log.WriteToStderrf("failed to persist compliance_operator_rules to store %v", err)
			return err
		}
	}
	return nil
}

func init() {
	migrations.MustRegisterMigration(migration)
}
