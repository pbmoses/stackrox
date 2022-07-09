package n39ton40

import (
	"context"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_39_to_n_40_postgres_policies/legacy"
	"github.com/stackrox/rox/pkg/bolthelper"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
	bolt "go.etcd.io/bbolt"
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
	legacyDB *bolt.DB

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
	s.legacyDB, err = bolthelper.NewTemp(s.T().Name() + ".db")
	s.NoError(err)

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)

	s.ctx = sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowAllAccessScopeChecker())
	s.pool, err = pgxpool.ConnectConfig(s.ctx, config)
	s.Require().NoError(err)
	pgtest.CleanUpDB(s.ctx, s.T(), s.pool)
	s.gormDB = pgtest.OpenGormDB(s.T(), source)
}

func (s *postgresMigrationSuite) TearDownTest() {
	testutils.TearDownDB(s.legacyDB)
	_ = s.gormDB.Migrator().DropTable(pkgSchema.CreateTablePoliciesStmt.GormModel)
	pgtest.CleanUpDB(s.ctx, s.T(), s.pool)
	pgtest.CloseGormDB(s.T(), s.gormDB)
	s.pool.Close()
}
func (s *postgresMigrationSuite) TestMigration() {
	// Prepare data and write to legacy DB
	var policys []*storage.Policy
	legacyStore := legacy.New(s.legacyDB)
	for i := 0; i < 200; i++ {
		policy := &storage.Policy{}
		s.NoError(testutils.FullInit(policy, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		policys = append(policys, policy)
		s.NoError(legacyStore.Upsert(s.ctx, policy))
	}
	s.NoError(move(s.legacyDB, s.gormDB, s.pool, legacyStore))
	var count int64
	s.gormDB.Model(pkgSchema.CreateTablePoliciesStmt.GormModel).Count(&count)
	s.Equal(int64(len(policys)), count)
	for _, policy := range policys {
		s.Equal(policy, s.get(policy.GetId()))
	}
}

func (s *postgresMigrationSuite) get(id string) *storage.Policy {

	q := search.ConjunctionQuery(
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	data, err := postgres.RunGetQueryForSchema(s.ctx, schema, q, s.pool)
	s.NoError(err)
	var msg storage.Policy
	s.NoError(proto.Unmarshal(data, &msg))
	return &msg
}
