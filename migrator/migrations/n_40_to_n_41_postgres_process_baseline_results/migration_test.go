// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package n40ton41

import (
	"context"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_40_to_n_41_postgres_process_baseline_results/legacy"

	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stackrox/rox/pkg/testutils/rocksdbtest"
	"github.com/stretchr/testify/suite"
	"github.com/tecbot/gorocksdb"
	"gorm.io/gorm"
)

func TestMigration(t *testing.T) {
	suite.Run(t, new(postgresMigrationSuite))
}

type postgresMigrationSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	ctx         context.Context

	// LegacyDB to migrate from
	legacyDB *rocksdb.RocksDB

	// PostgresDB
	pool   *pgxpool.Pool
	gormDB *gorm.DB
}

var _ suite.TearDownTestSuite = (*postgresMigrationSuite)(nil)

func (s *postgresMigrationSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")
	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	var err error
	s.legacyDB, err = rocksdb.NewTemp(s.T().Name())
	s.NoError(err)

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)

	s.ctx = sac.WithAllAccess(context.Background())
	s.pool, err = pgxpool.ConnectConfig(s.ctx, config)
	s.Require().NoError(err)
	pgtest.CleanUpDB(s.ctx, s.T(), s.pool)
	s.gormDB = pgtest.OpenGormDBWithDisabledConstraints(s.T(), source)
}

func (s *postgresMigrationSuite) TearDownTest() {
	rocksdbtest.TearDownRocksDB(s.legacyDB)
	_ = s.gormDB.Migrator().DropTable(pkgSchema.CreateTableProcessBaselineResultsStmt.GormModel)
	pgtest.CleanUpDB(s.ctx, s.T(), s.pool)
	pgtest.CloseGormDB(s.T(), s.gormDB)
	s.pool.Close()
}
func (s *postgresMigrationSuite) TestMigration() {
	// Prepare data and write to legacy DB
	var processBaselineResultss []*storage.ProcessBaselineResults
	legacyStore, err := legacy.New(s.legacyDB)
	s.NoError(err)
	batchSize = 48
	rocksWriteBatch := gorocksdb.NewWriteBatch()
	defer rocksWriteBatch.Destroy()
	for i := 0; i < 200; i++ {
		processBaselineResults := &storage.ProcessBaselineResults{}
		s.NoError(testutils.FullInit(processBaselineResults, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		processBaselineResultss = append(processBaselineResultss, processBaselineResults)
	}
	s.NoError(legacyStore.UpsertMany(s.ctx, processBaselineResultss))
	s.NoError(move(s.gormDB, s.pool, legacyStore))
	var count int64
	s.gormDB.Model(pkgSchema.CreateTableProcessBaselineResultsStmt.GormModel).Count(&count)
	s.Equal(int64(len(processBaselineResultss)), count)
	for _, processBaselineResults := range processBaselineResultss {
		s.Equal(processBaselineResults, s.get(processBaselineResults.GetDeploymentId()))
	}
}

func (s *postgresMigrationSuite) get(deploymentId string) *storage.ProcessBaselineResults {

	q := search.ConjunctionQuery(
		search.NewQueryBuilder().AddDocIDs(deploymentId).ProtoQuery(),
	)

	data, err := postgres.RunGetQueryForSchema(s.ctx, schema, q, s.pool)
	s.NoError(err)
	var msg storage.ProcessBaselineResults
	s.NoError(proto.Unmarshal(data, &msg))
	return &msg
}
