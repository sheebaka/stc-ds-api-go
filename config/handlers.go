package config

import (
	"context"
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"github.com/gin-gonic/gin"
)

type SfAccount struct {
	ID               string `json:"ID"`
	Name             string `json:"Name"`
	DOTNumber        string `json:"DOT_Number__c"`
	CRMAccountNumber string `json:"CRM_Account_Number__c"`
}

func (a *AppConfig) Account(c *gin.Context) {
	dotNumber := c.Query("dotnumber")
	cols := core.StringSlice{"Id", "Name", "DOT_Number__c", "CRM_Account_Number__c"}
	query := fmt.Sprintf(`SELECT %s FROM dev_silver.crm.sf_account WHERE DOT_Number__c=%s`, cols.Join(", "), dotNumber)
	rows, err := a.SqlDB.QueryContext(context.Background(), query)
	if err != nil {
		err = fmt.Errorf("error querying accounts: %s", err)
		fmt.Println(err)
		return
	}
	defer rows.Close()
	account := SfAccount{}
	for rows.Next() {
		err = rows.Scan(&account.ID, &account.Name, &account.DOTNumber, &account.CRMAccountNumber)
		if err != nil {
			return
		}
		c.JSON(200, account)
	}
	//dotNumber := c.Query("dotnumber")
	//if dotNumber == "" {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "dotnumber is required"})
	//	return
	//}
	//q := query..Use(a.GormDB)
	////account := model.SfAccount{}
	////stmt := fmt.Sprintf("SELECT * FROM sf_account WHERE DOT_Number__c = %s", dotNumber)
	//expr := q.SfAccount.DOTNumberC.Eq(dotNumber)
	//account, err := q.SfAccount.Where(expr).First()
	////tx := q.SfAccount.UnderlyingDB().Raw(stmt).First(&account)
	//if err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	//	return
	//}
	//c.JSON(http.StatusOK, gin.H{"account": account})
}

// ========================
