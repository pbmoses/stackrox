package datastore

import (
	"context"
	"testing"

	"github.com/blevesearch/bleve"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/central/globalindex"
	"github.com/stackrox/rox/central/processbaseline/index"
	"github.com/stackrox/rox/central/processbaseline/index/mappings"
	"github.com/stackrox/rox/central/processbaseline/search"
	"github.com/stackrox/rox/central/processbaseline/store"
	pgStore "github.com/stackrox/rox/central/processbaseline/store/postgres"
	rdbStore "github.com/stackrox/rox/central/processbaseline/store/rocksdb"
	processBaselineResultMock "github.com/stackrox/rox/central/processbaselineresults/datastore/mocks"
	processIndicatorMock "github.com/stackrox/rox/central/processindicator/datastore/mocks"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sac/testconsts"
	"github.com/stackrox/rox/pkg/sac/testutils"
	searchPkg "github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/suite"
)

func TestProcessBaselineDatastoreSAC(t *testing.T) {
	suite.Run(t, new(processBaselineSACTestSuite))
}

type processBaselineSACTestSuite struct {
	suite.Suite

	engine *rocksdb.RocksDB
	index  bleve.Index

	pool *pgxpool.Pool

	storage store.Store
	indexer index.Indexer
	search  search.Searcher

	datastore DataStore

	optionsMap searchPkg.OptionsMap

	testContexts map[string]context.Context

	processBaselineResultMock *processBaselineResultMock.MockDataStore
	processIndicatorMock      *processIndicatorMock.MockDataStore

	testProcessBaselineIDs []string
}

func (s *processBaselineSACTestSuite) SetupSuite() {
	var err error
	if features.PostgresDatastore.Enabled() {
		ctx := context.Background()
		src := pgtest.GetConnectionString(s.T())
		cfg, err := pgxpool.ParseConfig(src)
		s.Require().NoError(err)
		s.pool, err = pgxpool.ConnectConfig(ctx, cfg)
		s.Require().NoError(err)
		pgStore.Destroy(ctx, s.pool)
		gormDB := pgtest.OpenGormDB(s.T(), src)
		defer pgtest.CloseGormDB(s.T(), gormDB)
		s.storage = pgStore.CreateTableAndNewStore(ctx, s.pool, gormDB)
		s.indexer = pgStore.NewIndexer(s.pool)
		s.optionsMap = schema.ProcessBaselinesSchema.OptionsMap
	} else {
		s.engine, err = rocksdb.NewTemp("processBaselineSACTest")
		s.Require().NoError(err)
		bleveIndex, err := globalindex.MemOnlyIndex()
		s.Require().NoError(err)
		s.index = bleveIndex

		s.storage, err = rdbStore.New(s.engine)
		s.Require().NoError(err)
		s.indexer = index.New(s.index)
		s.optionsMap = mappings.OptionsMap
	}

	s.search, err = search.New(s.storage, s.indexer)
	s.Require().NoError(err)

	s.processIndicatorMock = processIndicatorMock.NewMockDataStore(gomock.NewController(s.T()))
	s.processBaselineResultMock = processBaselineResultMock.NewMockDataStore(gomock.NewController(s.T()))

	s.datastore = New(s.storage, s.indexer, s.search, s.processBaselineResultMock, s.processIndicatorMock)

	s.testContexts = testutils.GetNamespaceScopedTestContexts(context.Background(), s.T(),
		resources.ProcessWhitelist)
}

func (s *processBaselineSACTestSuite) TearDownSuite() {
	if features.PostgresDatastore.Enabled() {
		s.pool.Close()
	} else {
		s.Require().NoError(rocksdb.CloseAndRemove(s.engine))
		s.Require().NoError(s.index.Close())
	}
}

func (s *processBaselineSACTestSuite) SetupTest() {
	s.testProcessBaselineIDs = make([]string, 0)

	processBaselines := fixtures.GetSACTestStorageProcessBaselineSet(fixtures.GetScopedProcessBaseline)

	for i := range processBaselines {
		_, err := s.datastore.AddProcessBaseline(s.testContexts[testutils.UnrestrictedReadWriteCtx],
			processBaselines[i])
		s.Require().NoError(err)
	}

	for _, rb := range processBaselines {
		s.testProcessBaselineIDs = append(s.testProcessBaselineIDs, rb.GetId())
	}
}

func (s *processBaselineSACTestSuite) TearDownTest() {
	s.Require().NoError(s.datastore.RemoveProcessBaselinesByIDs(s.testContexts[testutils.UnrestrictedReadWriteCtx],
		s.testProcessBaselineIDs))
}

func (s *processBaselineSACTestSuite) deleteProcessBaseline(id string) {
	if id != "" {
		s.Require().NoError(s.datastore.RemoveProcessBaselinesByIDs(s.testContexts[testutils.UnrestrictedReadWriteCtx],
			[]string{id}))
	}
}

func (s *processBaselineSACTestSuite) TestAddProcessBaseline() {
	testedVerb := "add"
	cases := testutils.GenericNamespaceSACUpsertTestCases(s.T(), testedVerb)

	for name, c := range cases {
		s.Run(name, func() {
			processBaseline := fixtures.GetScopedProcessBaseline(uuid.NewV4().String(), testconsts.Cluster2,
				testconsts.NamespaceB)
			ctx := s.testContexts[c.ScopeKey]
			key, err := s.datastore.AddProcessBaseline(ctx, processBaseline)
			defer s.deleteProcessBaseline(key)
			if c.ExpectError {
				s.Require().Error(err)
				s.ErrorIs(err, c.ExpectedError)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *processBaselineSACTestSuite) TestUpsertProcessBaseline() {
	testedVerb := "upsert"
	cases := testutils.GenericNamespaceSACUpsertTestCases(s.T(), testedVerb)

	for name, c := range cases {
		s.Run(name, func() {
			processBaseline := fixtures.GetScopedProcessBaseline(uuid.NewV4().String(), testconsts.Cluster2,
				testconsts.NamespaceB)
			ctx := s.testContexts[c.ScopeKey]
			processBaseline, err := s.datastore.UpsertProcessBaseline(ctx, processBaseline.GetKey(), nil, false, false)
			defer s.deleteProcessBaseline(processBaseline.GetId())
			if c.ExpectError {
				s.Require().Error(err)
				s.ErrorIs(err, c.ExpectedError)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *processBaselineSACTestSuite) TestUpdateProcessBaselineElements() {
	processBaseline := fixtures.GetScopedProcessBaseline(uuid.NewV4().String(), testconsts.Cluster2,
		testconsts.NamespaceB)
	_, err := s.datastore.AddProcessBaseline(s.testContexts[testutils.UnrestrictedReadWriteCtx], processBaseline)
	s.Require().NoError(err)
	s.testProcessBaselineIDs = append(s.testProcessBaselineIDs, processBaseline.GetId())

	testedVerb := "update"
	cases := testutils.GenericNamespaceSACUpsertTestCases(s.T(), testedVerb)

	for name, c := range cases {
		s.Run(name, func() {
			ctx := s.testContexts[c.ScopeKey]
			_, err := s.datastore.UpdateProcessBaselineElements(
				ctx, processBaseline.GetKey(), nil, nil, false)
			if c.ExpectError {
				s.Require().Error(err)
				s.ErrorIs(err, c.ExpectedError)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *processBaselineSACTestSuite) TestGetProcessBaseline() {
	processBaseline := fixtures.GetScopedProcessBaseline(uuid.NewV4().String(), testconsts.Cluster2,
		testconsts.NamespaceB)
	_, err := s.datastore.AddProcessBaseline(s.testContexts[testutils.UnrestrictedReadWriteCtx], processBaseline)
	s.Require().NoError(err)
	s.testProcessBaselineIDs = append(s.testProcessBaselineIDs, processBaseline.GetId())

	cases := testutils.GenericNamespaceSACGetTestCases(s.T())

	for name, c := range cases {
		s.Run(name, func() {
			ctx := s.testContexts[c.ScopeKey]
			res, found, err := s.datastore.GetProcessBaseline(ctx, processBaseline.GetKey())
			s.Require().NoError(err)
			if c.ExpectedFound {
				s.Require().True(found)
				s.Equal(*processBaseline, *res)
			} else {
				s.False(found)
				s.Nil(res)
			}
		})
	}
}

func (s *processBaselineSACTestSuite) TestRemoveProcessBaseline() {
	cases := testutils.GenericNamespaceSACDeleteTestCases(s.T())

	s.processBaselineResultMock.EXPECT().DeleteBaselineResults(gomock.Any(), gomock.Any()).Return(nil).
		AnyTimes()

	for name, c := range cases {
		s.Run(name, func() {
			processBaseline := fixtures.GetScopedProcessBaseline(uuid.NewV4().String(), testconsts.Cluster2,
				testconsts.NamespaceB)
			_, err := s.datastore.AddProcessBaseline(s.testContexts[testutils.UnrestrictedReadWriteCtx], processBaseline)
			s.Require().NoError(err)
			s.testProcessBaselineIDs = append(s.testProcessBaselineIDs, processBaseline.GetId())
			defer s.deleteProcessBaseline(processBaseline.GetId())

			ctx := s.testContexts[c.ScopeKey]
			err = s.datastore.RemoveProcessBaseline(ctx, processBaseline.GetKey())
			if c.ExpectError {
				s.Require().Error(err)
				s.ErrorIs(err, c.ExpectedError)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *processBaselineSACTestSuite) runSearchRawTest(c testutils.SACSearchTestCase) {
	ctx := s.testContexts[c.ScopeKey]
	results, err := s.datastore.SearchRawProcessBaselines(ctx, nil)
	s.Require().NoError(err)
	resultObjs := make([]sac.NamespaceScopedObject, 0, len(results))
	for i := range results {
		resultObjs = append(resultObjs, results[i].GetKey())
	}
	resultCounts := testutils.CountSearchResultObjectsPerClusterAndNamespace(s.T(), resultObjs)
	testutils.ValidateSACSearchResultDistribution(&s.Suite, c.Results, resultCounts)
}

func (s *processBaselineSACTestSuite) runSearchTest(c testutils.SACSearchTestCase) {
	ctx := s.testContexts[c.ScopeKey]
	results, err := s.datastore.Search(ctx, nil)
	s.Require().NoError(err)
	resultCounts := testutils.CountResultsPerClusterAndNamespace(s.T(), results, s.optionsMap)
	testutils.ValidateSACSearchResultDistribution(&s.Suite, c.Results, resultCounts)
}

func (s *processBaselineSACTestSuite) TestScopedSearch() {
	for name, c := range testutils.GenericScopedSACSearchTestCases(s.T()) {
		s.Run(name, func() {
			s.runSearchTest(c)
		})
	}
}

func (s *processBaselineSACTestSuite) TestUnrestrictedSearch() {
	for name, c := range testutils.GenericUnrestrictedSACSearchTestCases(s.T()) {
		s.Run(name, func() {
			s.runSearchTest(c)
		})
	}
}

func (s *processBaselineSACTestSuite) TestScopeSearchRaw() {
	for name, c := range testutils.GenericScopedSACSearchTestCases(s.T()) {
		s.Run(name, func() {
			s.runSearchRawTest(c)
		})
	}
}

func (s *processBaselineSACTestSuite) TestUnrestrictedSearchRaw() {
	for name, c := range testutils.GenericUnrestrictedRawSACSearchTestCases(s.T()) {
		s.Run(name, func() {
			s.runSearchRawTest(c)
		})
	}
}
