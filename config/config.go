package config

import (
	"database/sql"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
)

// ========================

type AppConfig struct {
	DataSources   `yaml:"sources"`
	*SourceConfig `yaml:"-"`
	CommonModel   Model `yaml:"model"`
}

type DataSources struct {
	core.Map[*SourceConfig] `yaml:",inline"`
}

type SourceConfig struct {
	Env            string `yaml:"env" yaml:"envr"`
	DriverConfig   `yaml:"config"`
	DataSourceName string `yaml:"dsn"`
	Model          `yaml:"model"`
	GormDB         *gorm.DB
	SqlDB          *sql.DB
	migrator.Migrator
	gorm.Dialector
}

type Model struct {
	ModelOutPath string           `yaml:"outpath"`
	ModelFile    string           `yaml:"filename"`
	ModelPkgPath string           `yaml:"pkgpath"`
	TableNames   core.StringSlice `yaml:"tables"`
	//ModelTableName string `yaml:"tablename"`
}

type DriverConfig struct {
	Database    string `yaml:"database" json:"dbname"`
	Host        string `yaml:"host" json:"host" env:"HOST"`
	Password    string `yaml:"password" json:"password" env:"DATABRICKS_TOKEN"`
	User        string `yaml:"user" json:"username" env:"DATABRICKS_USER"`
	Port        string `yaml:"port" json:"port"`
	Path        string `yaml:"path"`
	Schema      string `yaml:"schema"`
	WarehouseId string `yaml:"warehouse_id" env:"DATABRICKS_WAREHOUSE_ID"`
	Secret      string `yaml:"secret"`
	Driver      string `yaml:"-"`
}
