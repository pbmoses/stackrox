// Code generated by pg-bindings generator. DO NOT EDIT.

package n4ton5

import (
	"context"
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_04_to_n_05_postgres_images/legacy"
	psImageStore "github.com/stackrox/rox/migrator/migrations/n_04_to_n_05_postgres_images/postgres"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stackrox/rox/pkg/testutils/rocksdbtest"
	"github.com/stretchr/testify/suite"
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
	_ = s.gormDB.Migrator().DropTable(pkgSchema.CreateTableImagesStmt.GormModel)
	pgtest.CleanUpDB(s.ctx, s.T(), s.pool)
	pgtest.CloseGormDB(s.T(), s.gormDB)
	s.pool.Close()
}
func (s *postgresMigrationSuite) TestMigration() {
	dacky, err := dackbox.NewRocksDBDackBox(s.legacyDB, nil, []byte("graph"), []byte("dirty"), []byte("valid"))
	s.NoError(err)
	legacyStore := legacy.New(dacky, concurrency.NewKeyFence(), false)
	s.NoError(err)
	// ts := &types.Timestamp{Seconds: 1245367890}

	images := []*storage.Image{
		{
			Id: "sha256:sha1",
			Name: &storage.ImageName{
				FullName: "name1",
			},
			Metadata: &storage.ImageMetadata{
				V1: &storage.V1Metadata{
					Created: types.TimestampNow(),
				},
			},
			Scan: &storage.ImageScan{
				ScanTime: types.TimestampNow(),
				Components: []*storage.EmbeddedImageScanComponent{
					{
						Name:    "comp1",
						Version: "ver1",
						HasLayerIndex: &storage.EmbeddedImageScanComponent_LayerIndex{
							LayerIndex: 1,
						},
						Vulns: []*storage.EmbeddedVulnerability{},
					},
					{
						Name:    "comp1",
						Version: "ver2",
						HasLayerIndex: &storage.EmbeddedImageScanComponent_LayerIndex{
							LayerIndex: 3,
						},
						Vulns: []*storage.EmbeddedVulnerability{
							{
								Cve:               "cve1",
								VulnerabilityType: storage.EmbeddedVulnerability_IMAGE_VULNERABILITY,
							},
							{
								Cve:               "cve2",
								VulnerabilityType: storage.EmbeddedVulnerability_IMAGE_VULNERABILITY,
								SetFixedBy: &storage.EmbeddedVulnerability_FixedBy{
									FixedBy: "ver3",
								},
							},
						},
					},
					{
						Name:    "comp2",
						Version: "ver1",
						HasLayerIndex: &storage.EmbeddedImageScanComponent_LayerIndex{
							LayerIndex: 2,
						},
						Vulns: []*storage.EmbeddedVulnerability{
							{
								Cve:               "cve1",
								VulnerabilityType: storage.EmbeddedVulnerability_IMAGE_VULNERABILITY,
								SetFixedBy: &storage.EmbeddedVulnerability_FixedBy{
									FixedBy: "ver2",
								},
							},
							{
								Cve:               "cve2",
								VulnerabilityType: storage.EmbeddedVulnerability_IMAGE_VULNERABILITY,
							},
						},
					},
				},
			},
			RiskScore: 30,
		},
		{
			Id: "sha256:sha2",
			Name: &storage.ImageName{
				FullName: "name2",
			},
		},
	}

	for _, image := range images {
		s.NoError(legacyStore.Upsert(s.ctx, image))
	}
	s.NoError(move(s.gormDB, s.pool, legacyStore))

	postgresStore := psImageStore.New(s.pool, true)
	for _, image := range images {
		got, found, err := postgresStore.Get(s.ctx, image.GetId())
		s.NoError(err)
		s.True(found)
		var gotComponent *storage.EmbeddedImageScanComponent
		for ci, component := range image.GetScan().GetComponents() {
			if len(got.GetScan().GetComponents()) > ci {
				gotComponent = got.GetScan().GetComponents()[ci]
			}
			for vi, vuln := range component.GetVulns() {
				vuln.VulnerabilityType = storage.EmbeddedVulnerability_IMAGE_VULNERABILITY
				vuln.VulnerabilityTypes = []storage.EmbeddedVulnerability_VulnerabilityType{storage.EmbeddedVulnerability_IMAGE_VULNERABILITY}
				if gotComponent == nil || len(gotComponent.GetVulns()) <= vi {
					continue
				}
				vuln.FirstSystemOccurrence = got.GetScan().GetComponents()[ci].GetVulns()[vi].GetFirstSystemOccurrence()
				vuln.FirstImageOccurrence = got.GetScan().GetComponents()[ci].GetVulns()[vi].GetFirstSystemOccurrence()
			}
		}
		s.Equal(image, got)

		listGot, exists, err := legacyStore.GetImageMetadata(s.ctx, image.GetId())
		s.NoError(err)
		s.True(exists)
		s.Equal(image.GetName().GetFullName(), listGot.GetName().GetFullName())
	}

	// Test Count
	count, err := postgresStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(images), count)
}

