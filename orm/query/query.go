package query

import (
	"github.com/stc-ds-databricks-go/orm/model"
	"gorm.io/gorm"
)

func FilterWithColumn(db *gorm.DB, col string, val string) (out model.SfAccount, err error) {
	out, err = Use(db).SfAccount.FilterWithColumn(col, val)
	return
}

const FilterWithColumnTemplate = `

func {{.FilterMethodName}}(db *gorm.DB, col string, val string) (out model.{{.ModelName}}, err error) {
	out, err = Use(db).{{.ModelName}}.{{.FilterMethodName}}(col, val)
}

`
