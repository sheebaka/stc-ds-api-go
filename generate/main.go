package main

import (
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"github.com/stc-ds-databricks-go/config"
	"github.com/stc-ds-databricks-go/generate/common"
	"gorm.io/gen"
)

const DriverName = config.Databricks

func GenerateModel(a *config.AppConfig) (err error) {
	g := gen.NewGenerator(gen.Config{
		OutPath:      config.JoinRoot(a.ModelOutPath),
		ModelPkgPath: config.JoinRoot(a.ModelPkgPath),
		Mode:         gen.WithoutContext | gen.WithQueryInterface | gen.WithDefaultQuery,
	})
	a.GormDB.Dialector = config.NewDatabricksDialector(a)
	g.UseDB(a.GormDB)
	//
	models := core.NewSlice[any]()
	for _, table := range a.TableNames {
		model := g.GenerateModelAs(table, config.ToTitleCase(table), gen.WithMethod(common.SharedMethods{}))
		models = models.AppendPtr(model)
	}
	//
	g.ApplyInterface(
		func(common.Filter) {},
		*models...,
	//g.GenerateModelAs(a.Table(), a.ModelName(), gen.WithMethod(common.CommonMethod{})),
	)
	g.Execute()
	return
}

func main() {
	app, err := config.ConfigureApp(DriverName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = GenerateModel(app); err != nil {
		fmt.Println(err)
		return
	}
}
