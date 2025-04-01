package model

const FilterWithColumnTemplate = `

func {{.FilterMethodName}}(db *gorm.DB, col string, val string) (out model.{{.ModelName}}, err error) {
	out, err = Use(db).{{.ModelName}}.{{.FilterMethodName}}(col, val)
}

`
