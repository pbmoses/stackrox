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

type NetworkBaselinesStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	store       Store
	pool        *pgxpool.Pool
}

func TestNetworkBaselinesStore(t *testing.T) {
	suite.Run(t, new(NetworkBaselinesStoreSuite))
}

func (s *NetworkBaselinesStoreSuite) SetupSuite() {
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
	gormDB := pgtest.OpenGormDB(s.T(), source, false)
	defer pgtest.CloseGormDB(s.T(), gormDB)
	s.store = CreateTableAndNewStore(ctx, pool, gormDB)
}

func (s *NetworkBaselinesStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.pool.Exec(ctx, "TRUNCATE network_baselines CASCADE")
	s.T().Log("network_baselines", tag)
	s.NoError(err)
}

func (s *NetworkBaselinesStoreSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
	s.envIsolator.RestoreAll()
}

func (s *NetworkBaselinesStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	networkBaseline := &storage.NetworkBaseline{}
	s.NoError(testutils.FullInit(networkBaseline, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundNetworkBaseline, exists, err := store.Get(ctx, networkBaseline.GetDeploymentId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNetworkBaseline)

	s.NoError(store.Upsert(ctx, networkBaseline))
	foundNetworkBaseline, exists, err = store.Get(ctx, networkBaseline.GetDeploymentId())
	s.NoError(err)
	s.True(exists)
	s.Equal(networkBaseline, foundNetworkBaseline)

	networkBaselineCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, networkBaselineCount)

	networkBaselineExists, err := store.Exists(ctx, networkBaseline.GetDeploymentId())
	s.NoError(err)
	s.True(networkBaselineExists)
	s.NoError(store.Upsert(ctx, networkBaseline))

	foundNetworkBaseline, exists, err = store.Get(ctx, networkBaseline.GetDeploymentId())
	s.NoError(err)
	s.True(exists)
	s.Equal(networkBaseline, foundNetworkBaseline)

	s.NoError(store.Delete(ctx, networkBaseline.GetDeploymentId()))
	foundNetworkBaseline, exists, err = store.Get(ctx, networkBaseline.GetDeploymentId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNetworkBaseline)

	var networkBaselines []*storage.NetworkBaseline
	for i := 0; i < 200; i++ {
		networkBaseline := &storage.NetworkBaseline{}
		s.NoError(testutils.FullInit(networkBaseline, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		networkBaselines = append(networkBaselines, networkBaseline)
	}

	s.NoError(store.UpsertMany(ctx, networkBaselines))

	networkBaselineCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, networkBaselineCount)
}
