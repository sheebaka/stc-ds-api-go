package orm

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
	//// SELECT * FROM @@table WHERE "DOT_Number__c"=@value
	//FilterByDotNumber(value string) (gen.T, error)
}

type CommonMethod struct {
	ID         string
	DOTNumberC string
}

func (m CommonMethod) GetId() string {
	return m.ID
}

func NewCommonMethod(id, dotNumberC string) *CommonMethod {
	return &CommonMethod{}
}

type Identifier interface {
	GetId() string
}
