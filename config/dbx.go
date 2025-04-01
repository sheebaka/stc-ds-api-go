package config

import (
	"database/sql"
	dbsql "github.com/databricks/databricks-sql-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
	"log"
	"strconv"
)

func (x *SourceConfig) newConnectorDB() {
	connector, err := dbsql.NewConnector(
		dbsql.WithServerHostname(x.Host),
		dbsql.WithPort(x.PortInt()),
		dbsql.WithHTTPPath(x.Path),
		dbsql.WithAccessToken(x.Password),
		dbsql.WithInitialNamespace(x.Database, x.Schema),
	)
	if err != nil {
		log.Fatal(err)
	}
	x.SqlDB = sql.OpenDB(connector)
}

func (d *DriverConfig) PortInt() (out int) {
	port, err := strconv.Atoi(d.Port)
	if err != nil {
		return
	}
	out = port
	return
}

func NewDatabricksDialector(app *AppConfig) DatabricksDialector {
	return DatabricksDialector{app, app.GormDB}
}

type DatabricksDialector struct {
	*AppConfig
	*gorm.DB
}

func (d DatabricksDialector) Name() string {
	return "databricks"
}

func (d DatabricksDialector) Initialize(db *gorm.DB) error {
	return d.Use(db)
}

type DatabricksMigrator struct {
	*migrator.Migrator
	*migrator.Config
	*AppConfig
}

func (d DatabricksDialector) Migrator(db *gorm.DB) gorm.Migrator {
	m := migrator.Migrator{
		Config: migrator.Config{
			CreateIndexAfterCreateTable: false,
			DB:                          db,
			Dialector:                   d.Dialector(),
		},
	}
	return &DatabricksMigrator{
		Migrator:  &m,
		Config:    &m.Config,
		AppConfig: d.AppConfig,
	}
}

func (d DatabricksDialector) DataTypeOf(field *schema.Field) string {
	// Implement data type mapping
	return field.FieldType.String()
}

func (d DatabricksDialector) DefaultValueOf(field *schema.Field) clause.Expression {
	// Implement default value handling
	return d.GormDB.DefaultValueOf(field)
}

func (d DatabricksDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, value interface{}) {
	writer.WriteByte('?')
}

func (d DatabricksDialector) QuoteTo(writer clause.Writer, fieldName string) {
	d.AppConfig.Dialector.QuoteTo(writer, fieldName)
}

func (d DatabricksDialector) Explain(sql string, vars ...interface{}) string {
	return d.AppConfig.Dialector.Explain(sql, vars...)
}

func (d DatabricksDialector) PrepareStmt(stmt *gorm.Statement) error {
	// Implement prepare statement logic
	return stmt.Use(d.GormDB)
}

func (d DatabricksDialector) ConnPool() gorm.ConnPool {
	// Implement connection pool logic
	return nil
}

func (d DatabricksDialector) Dialector() DatabricksDialector {
	return d
}
