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

type TestG3GrandChild1StoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	store       Store
	pool        *pgxpool.Pool
}

func TestTestG3GrandChild1Store(t *testing.T) {
	suite.Run(t, new(TestG3GrandChild1StoreSuite))
}

func (s *TestG3GrandChild1StoreSuite) SetupTest() {
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

func (s *TestG3GrandChild1StoreSuite) TearDownTest() {
	if s.pool != nil {
		s.pool.Close()
	}
	s.envIsolator.RestoreAll()
}

func (s *TestG3GrandChild1StoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	testG3GrandChild1 := &storage.TestG3GrandChild1{}
	s.NoError(testutils.FullInit(testG3GrandChild1, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundTestG3GrandChild1, exists, err := store.Get(ctx, testG3GrandChild1.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTestG3GrandChild1)

	s.NoError(store.Upsert(ctx, testG3GrandChild1))
	foundTestG3GrandChild1, exists, err = store.Get(ctx, testG3GrandChild1.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(testG3GrandChild1, foundTestG3GrandChild1)

	testG3GrandChild1Count, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, testG3GrandChild1Count)

	testG3GrandChild1Exists, err := store.Exists(ctx, testG3GrandChild1.GetId())
	s.NoError(err)
	s.True(testG3GrandChild1Exists)
	s.NoError(store.Upsert(ctx, testG3GrandChild1))

	foundTestG3GrandChild1, exists, err = store.Get(ctx, testG3GrandChild1.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(testG3GrandChild1, foundTestG3GrandChild1)

	s.NoError(store.Delete(ctx, testG3GrandChild1.GetId()))
	foundTestG3GrandChild1, exists, err = store.Get(ctx, testG3GrandChild1.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTestG3GrandChild1)

	var testG3GrandChild1s []*storage.TestG3GrandChild1
	for i := 0; i < 200; i++ {
		testG3GrandChild1 := &storage.TestG3GrandChild1{}
		s.NoError(testutils.FullInit(testG3GrandChild1, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		testG3GrandChild1s = append(testG3GrandChild1s, testG3GrandChild1)
	}

	s.NoError(store.UpsertMany(ctx, testG3GrandChild1s))

	testG3GrandChild1Count, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, testG3GrandChild1Count)
}
