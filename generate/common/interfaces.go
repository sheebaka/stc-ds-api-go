package common

import (
	"gorm.io/gen"
)

// Querier Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

type Filter interface {
	// SELECT * FROM @@table WHERE @@column=@value
	FilterWithColumn(column string, value string) (gen.T, error)
}

//type FilterTempl interface {
//	// SELECT * FROM @@table WHERE @@column=@value
//	FilterWithColumn(column string, value string) (Identifier, error)
//}

type Identifier interface {
	GetId() string
	GetDOTNumber() string
	GetResponse() map[string]interface{}
}

//
//g.ApplyInterface(
//	func(orm.Filter) {},
//	g.GenerateAllTable(gen.WithMethod(orm.CommonMethod{}))...,
//)
//
