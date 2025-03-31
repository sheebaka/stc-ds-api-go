package config

import (
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"github.com/stc-ds-databricks-go/orm"
	"gorm.io/gen"
)

// ========================

func ConfigureApp(driverName string) (app *AppConfig, err error) {
	err = core.ReadYamlFile(JoinRoot("config.yaml"), &app)
	if err != nil {
		err = fmt.Errorf("failed to read config: %w", err)
		return
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
			return
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
	g.UseDB(a.GormDB)
	//
	g.ApplyInterface(
		func(orm.Filter) {},
		g.GenerateModelAs(a.SchemaDotTable(), a.ModelName(), gen.WithMethod(orm.CommonMethod{})),
	)
	//
	//g.ApplyInterface(
	//	func(orm.Filter) {},
	//	g.GenerateAllTable(gen.WithMethod(orm.CommonMethod{}))...,
	//)
	//
	g.Execute()
	return
}
