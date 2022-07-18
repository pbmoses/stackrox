// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableImageComponentsStmt holds the create statement for table `image_components`.
	CreateTableImageComponentsStmt = &postgres.CreateStmts{
		Table: `
               create table if not exists image_components (
                   Id varchar,
                   Name varchar,
                   Version varchar,
                   Priority integer,
                   Source integer,
                   RiskScore numeric,
                   TopCvss numeric,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		GormModel: (*ImageComponents)(nil),
		Indexes:   []string{},
		Children:  []*postgres.CreateStmts{},
	}

	// ImageComponentsSchema is the go schema for table `image_components`.
	ImageComponentsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("image_components")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ImageComponent)(nil)), "image_components")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_IMAGE_COMPONENTS, "imagecomponent", (*storage.ImageComponent)(nil)))
		schema.SetSearchScope([]v1.SearchCategory{
			v1.SearchCategory_IMAGE_VULNERABILITIES,
			v1.SearchCategory_COMPONENT_VULN_EDGE,
			v1.SearchCategory_IMAGE_COMPONENTS,
			v1.SearchCategory_IMAGE_COMPONENT_EDGE,
			v1.SearchCategory_IMAGE_VULN_EDGE,
			v1.SearchCategory_IMAGES,
			v1.SearchCategory_DEPLOYMENTS,
			v1.SearchCategory_NAMESPACES,
			v1.SearchCategory_CLUSTERS,
		}...)
		RegisterTable(schema, CreateTableImageComponentsStmt)
		return schema
	}()
)

const (
	ImageComponentsTableName = "image_components"
)

// ImageComponents holds the Gorm model for Postgres table `image_components`.
type ImageComponents struct {
	Id         string             `gorm:"column:id;type:varchar;primaryKey"`
	Name       string             `gorm:"column:name;type:varchar"`
	Version    string             `gorm:"column:version;type:varchar"`
	Priority   int64              `gorm:"column:priority;type:integer"`
	Source     storage.SourceType `gorm:"column:source;type:integer"`
	RiskScore  float32            `gorm:"column:riskscore;type:numeric"`
	TopCvss    float32            `gorm:"column:topcvss;type:numeric"`
	Serialized []byte             `gorm:"column:serialized;type:bytea"`
}
