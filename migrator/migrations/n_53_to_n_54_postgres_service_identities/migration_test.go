// Code generated by pg-bindings generator. DO NOT EDIT.

package n53ton54

import (
	"context"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_53_to_n_54_postgres_service_identities/legacy"
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
	pgtest.CleanUpDB(s.T(), s.ctx, s.pool)
	s.gormDB = pgtest.OpenGormDBWithDisabledConstraints(s.T(), source)
}

func (s *postgresMigrationSuite) TearDownTest() {
	testutils.TearDownDB(s.legacyDB)
	_ = s.gormDB.Migrator().DropTable(pkgSchema.CreateTableServiceIdentitiesStmt.GormModel)
	pgtest.CleanUpDB(s.T(), s.ctx, s.pool)
	pgtest.CloseGormDB(s.T(), s.gormDB)
	s.pool.Close()
}
func (s *postgresMigrationSuite) TestMigration() {
	// Prepare data and write to legacy DB
	var serviceIdentitys []*storage.ServiceIdentity
	legacyStore := legacy.New(s.legacyDB)
	for i := 0; i < 200; i++ {
		serviceIdentity := &storage.ServiceIdentity{}
		s.NoError(testutils.FullInit(serviceIdentity, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		serviceIdentitys = append(serviceIdentitys, serviceIdentity)
		s.NoError(legacyStore.Upsert(s.ctx, serviceIdentity))
	}
	s.NoError(move(s.legacyDB, s.gormDB, s.pool, legacyStore))
	var count int64
	s.gormDB.Model(pkgSchema.CreateTableServiceIdentitiesStmt.GormModel).Count(&count)
	s.Equal(int64(len(serviceIdentitys)), count)
	for _, serviceIdentity := range serviceIdentitys {
		s.Equal(serviceIdentity, s.get(serviceIdentity.GetSerialStr()))
	}
}

func (s *postgresMigrationSuite) get(serialStr string) *storage.ServiceIdentity {

	q := search.ConjunctionQuery(
		search.NewQueryBuilder().AddDocIDs(serialStr).ProtoQuery(),
	)

	data, err := postgres.RunGetQueryForSchema(s.ctx, schema, q, s.pool)
	s.NoError(err)
	var msg storage.ServiceIdentity
	s.NoError(proto.Unmarshal(data, &msg))
	return &msg
}
