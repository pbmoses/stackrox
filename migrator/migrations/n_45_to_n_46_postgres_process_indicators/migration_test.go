// Code generated by pg-bindings generator. DO NOT EDIT.

package n45ton46

import (
	"context"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_45_to_n_46_postgres_process_indicators/legacy"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/rocksdb"
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

	s.ctx = context.Background()
	s.pool, err = pgxpool.ConnectConfig(s.ctx, config)
	s.Require().NoError(err)
	pgtest.CleanUpDB(s.T(), s.ctx, s.pool)
	s.gormDB = pgtest.OpenGormDB(s.T(), source)
}

func (s *postgresMigrationSuite) TearDownTest() {
	rocksdbtest.TearDownRocksDB(s.legacyDB)
	_ = s.gormDB.Migrator().DropTable(pkgSchema.CreateTableProcessIndicatorsStmt.GormModel)
	pgtest.CleanUpDB(s.T(), s.ctx, s.pool)
	pgtest.CloseGormDB(s.T(), s.gormDB)
	s.pool.Close()
}

func (s *postgresMigrationSuite) TestMigration() {
	// Prepare data and write to legacy DB
	var processIndicators []*storage.ProcessIndicator
	legacyStore, err := legacy.New(s.legacyDB)
	s.NoError(err)
	batchSize = 48
	rocksWriteBatch := gorocksdb.NewWriteBatch()
	defer rocksWriteBatch.Destroy()
	for i := 0; i < 200; i++ {
		processIndicator := &storage.ProcessIndicator{}
		s.NoError(testutils.FullInit(processIndicator, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		processIndicators = append(processIndicators, processIndicator)
	}
	s.NoError(legacyStore.UpsertMany(s.ctx, processIndicators))
	s.NoError(moveProcessIndicators(s.legacyDB, s.gormDB, s.pool, legacyStore))
	var count int64
	s.gormDB.Model(pkgSchema.CreateTableProcessIndicatorsStmt.GormModel).Count(&count)
	s.Equal(int64(len(processIndicators)), count)
	for _, processIndicator := range processIndicators {
		s.Equal(processIndicator, s.get(processIndicator.GetId()))
	}
}

func (s *postgresMigrationSuite) get(id string) *storage.ProcessIndicator {

	q := search.ConjunctionQuery(
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	data, err := postgres.RunGetQueryForSchema(s.ctx, schema, q, s.pool)
	s.NoError(err)
	var msg storage.ProcessIndicator
	s.NoError(proto.Unmarshal(data, &msg))
	return &msg
}
