package schema

import (
	"strings"

	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/set"
	"gorm.io/gorm"
)

var (
	log = logging.LoggerForModule()
	// registeredTables is map of sql table name to go schema of the sql table.
	registeredTables = make(map[string]*RegisteredTable)
)

type RegisteredTable struct {
	Schema     *walker.Schema
	CreateStmt *postgres.CreateStmts
}

// RegisterTable maps a table to an object type for the purposes of metrics gathering
func RegisterTable(schema *walker.Schema, stmt *postgres.CreateStmts) {
	if _, ok := registeredTables[schema.Table]; ok {
		log.Fatalf("table %q is already registered for %s", schema.Table, schema.Type)
		return
	}
	registeredTables[schema.Table] = &RegisteredTable{Schema: schema, CreateStmt: stmt}
}

// GetSchemaForTable return the schema registered for specified table name.
func GetSchemaForTable(tableName string) *walker.Schema {
	if rt, ok := registeredTables[tableName]; ok {
		return rt.Schema
	}
	return nil
}

func getAllRegisteredTableInOrder() []*RegisteredTable {
	visited := set.NewStringSet()

	var rts []*RegisteredTable
	for table, _ := range registeredTables {
		rts = append(rts, getRegisteredTablesFor(visited, table)...)
	}
	return rts
}

func getRegisteredTablesFor(visited set.StringSet, table string) []*RegisteredTable {
	if visited.Contains(table) {
		return nil
	}
	var rts []*RegisteredTable
	rt := registeredTables[table]
	for _, ref := range rt.Schema.References {
		rts = append(rts, getRegisteredTablesFor(visited, ref.OtherSchema.Table)...)
	}
	rts = append(rts, rt)
	visited.Add(table)
	return rts
}

// ApplyCurrentSchema creates or auto migrate according to the current schema
func ApplyCurrentSchema(gormDB *gorm.DB) error {
	for _, rt := range getAllRegisteredTableInOrder() {
		// Exclude tests
		if strings.HasPrefix(rt.Schema.Table, "test_") {
			continue
		}
		pgutils.CreateTableFromModel(gormDB, rt.CreateStmt)
	}
	return nil
}
