package config

import (
	"github.com/stc-ds-databricks-go/orm/query/common"
	"gorm.io/gen"
)

func (a *AppConfig) GenerateModel() (err error) {
	g := gen.NewGenerator(gen.Config{
		OutPath:      JoinRoot(a.ModelOutPath),
		ModelPkgPath: JoinRoot(a.ModelPkgPath),
		Mode:         gen.WithoutContext | gen.WithQueryInterface | gen.WithDefaultQuery,
	})
	a.GormDB.Dialector = NewDatabricksDialector(a)
	g.UseDB(a.GormDB)
	g.ApplyInterface(
		func(common.Filter) {},
		g.GenerateModelAs(a.Table(), a.ModelName(), gen.WithMethod(common.CommonMethod{}.GetResponse)),
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
