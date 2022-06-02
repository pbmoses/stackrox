// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
)

type TestParent1StoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	store       Store
	pool        *pgxpool.Pool
}

func TestTestParent1Store(t *testing.T) {
	suite.Run(t, new(TestParent1StoreSuite))
}

func (s *TestParent1StoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	ctx := sac.WithAllAccess(context.Background())

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.Require().NoError(err)

	Destroy(ctx, pool)

	s.pool = pool
	gormDB := pgtest.OpenGormDB(s.T(), source)
	s.store = NewTestStore(ctx, pool, gormDB)
}

func (s *TestParent1StoreSuite) TearDownTest() {
	if s.pool != nil {
		s.pool.Close()
	}
	s.envIsolator.RestoreAll()
}

func (s *TestParent1StoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	testParent1 := &storage.TestParent1{}
	s.NoError(testutils.FullInit(testParent1, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundTestParent1, exists, err := store.Get(ctx, testParent1.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTestParent1)

}
