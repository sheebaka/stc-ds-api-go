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

type CommonMethod struct {
	Id                    string
	DOT_Number__c         string
	CRM_Account_Number__c string
	Name                  string
}

func (m CommonMethod) GetId() string {
	return m.Id
}

func (m CommonMethod) GetDOTNumber() string {
	return m.DOT_Number__c
}

func (m CommonMethod) GetResponse() (out map[string]interface{}) {
	out = make(map[string]interface{})
	out["Name"] = m.Name
	out["DOT_Number__c"] = m.DOT_Number__c
	out["CRM_Account_Number__c"] = m.CRM_Account_Number__c
	out["Id"] = m.Id
	return
}

type Identifier interface {
	GetId() string
	GetDOTNumber() string
	GetResponse() map[string]interface{}
}
