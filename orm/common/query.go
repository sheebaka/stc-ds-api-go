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

type CustomerStatus struct {
	Active__c string `json:"factoring,omitempty"`
}

type Result struct {
	Id                    *string `json:"-"`
	DOT_Number__c         int64   `json:"dotNumber"`
	CRM_Account_Number__c int     `json:"crmNumber"`
	Name                  string  `json:"companyName"`
	CustomerStatus        `json:"customerStatus,omitempty"`
}

func FilterByDotNumber(db *gorm.DB, dotNumber string) (result *Result, err error) {
	sfAccount := query.Use(db).SfAccount
	a := sfAccount.As("a")
	//
	err = a.Select(a.DOT_Number__c, a.CRM_Account_Number__c, a.Name, a.Id).Where(a.DOT_Number__c.Eq(dotNumber)).Scan(&result)

	sfCadenceDetails := query.Use(db).SfCadenceDetails
	c := sfCadenceDetails.As("c")
	//
	if result != nil {
		err = c.Select(c.Active__c).Where(c.Account__c.Eq(*result.Id), c.Active__c.Is(true)).Scan(&result.CustomerStatus)
	}
	//err = a.Select(a.DOT_Number__c, a.CRM_Account_Number__c, a.Name, c.Active__c).
	//	LeftJoin(c, a.Id.EqCol(c.Account__c)).
	//	Where(a.DOT_Number__c.Eq(dotNumber), c.Active__c.Is(true)).
	//	Scan(&results)
	return
}
