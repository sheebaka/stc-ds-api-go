package common

import (
	"github.com/stc-ds-databricks-go/orm/model"
	"github.com/stc-ds-databricks-go/orm/query"
	"gorm.io/gorm"
)

func FilterWithColumn(db *gorm.DB, col string, val string) (out model.SfAccount, err error) {
	out, err = query.Use(db).SfAccount.FilterWithColumn(col, val)
	return
}

// {
//"dotNumber": 3134772,
//"crmNumber": 9468470,
//"companyName": "SANGER EXPRESS INC",
//"customerStatus": {
//"factoring": "active"
//}
//}

type Result struct {
	DOT_Number__c         string `json:"dotNumber"`
	CRM_Account_Number__c string `json:"crmNumber"`
	Name                  string `json:"companyName"`
	Active__c             bool   `json:"active"`
}

func FilterByDotNumber(db *gorm.DB, dotNumber string) (results []*Result, err error) {
	sfAccount := query.Use(db).SfAccount
	sfCadenceDetails := query.Use(db).SfCadenceDetails
	a := sfAccount.As("a")
	c := sfCadenceDetails.As("c")
	err = c.Select(a.DOT_Number__c, a.CRM_Account_Number__c, a.Name, c.Active__c).LeftJoin(a, a.Id.EqCol(c.Account__c)).Where(a.DOT_Number__c.Eq(dotNumber)).Scan(&results)
	return
}

const FilterWithColumnTemplate = `

func {{.FilterMethodName}}(db *gorm.DB, col string, val string) (out model.{{.ModelName}}, err error) {
	out, err = Use(db).{{.ModelName}}.{{.FilterMethodName}}(col, val)
}

`
