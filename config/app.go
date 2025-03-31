package config

import (
	"database/sql"
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"github.com/stc-ds-databricks-go/orm"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// ========================

type AppConfig struct {
	DataSources   `yaml:"sources"`
	*SourceConfig `yaml:"-"`
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
}

type Model struct {
	ModelOutPath   string `yaml:"outpath"`
	ModelFile      string `yaml:"filename"`
	ModelPkgPath   string `yaml:"pkgpath"`
	ModelTableName string `yaml:"tablename"`
}

type DriverConfig struct {
	Database    string           `yaml:"database" json:"dbname"`
	Host        string           `yaml:"host" json:"host"`
	Password    string           `yaml:"password" json:"password"`
	User        string           `yaml:"user" json:"username"`
	Port        string           `yaml:"port" json:"port"`
	Path        string           `yaml:"path"`
	Schema      string           `yaml:"schema"`
	WarehouseId string           `yaml:"warehouse_id"`
	Secret      string           `yaml:"secret"`
	Driver      string           `yaml:"-"`
	Tables      core.StringSlice `yaml:"tables"`
}

// ========================

func ConfigureApp(driverName string) (app *AppConfig, err error) {
	err = core.ReadYamlFile(JoinRoot("config.yaml"), &app)
	if err != nil {
		err = fmt.Errorf("failed to read config: %w", err)
		panic(err)
	}
	sourceConfig, ok := app.FindAndGet(driverName)
	if ok && sourceConfig != nil {
		sourceConfig.Driver = driverName
		if sourceConfig.Secret != "" {
			if err = sourceConfig.getSecret(); err != nil {
				fmt.Println(err)
			}
		}
		if err = sourceConfig.resolve(); err != nil {
			err = fmt.Errorf("failed to configure database dsn: %v", err)
			panic(err)
		}
	}
	app.SourceConfig = sourceConfig
	if err = app.ConfigGormDB(); err != nil {
		err = fmt.Errorf("failed to get *gorm.DB: %w", err)
		return
	}
	return
}

func (a *AppConfig) GenerateModel() (err error) {
	//
	g := gen.NewGenerator(gen.Config{
		OutPath:      JoinRoot(a.ModelOutPath),
		ModelPkgPath: JoinRoot(a.ModelPkgPath),
		Mode:         gen.WithoutContext | gen.WithQueryInterface | gen.WithDefaultQuery, // generate mode
	})
	g.UseDB(a.GormDB) // reuse your gorm db
	//
	//meta := g.GenerateModel(a.TableName(), gen.WithMethod(orm.CommonMethod{}))
	//
	//g.ApplyInterface(func(orm.Filter) {}, meta)
	//
	g.ApplyInterface(
		func(orm.Filter) {},
		g.GenerateAllTable(gen.WithMethod(orm.CommonMethod{}))...,
	)
	//
	g.Execute()
	return
}
