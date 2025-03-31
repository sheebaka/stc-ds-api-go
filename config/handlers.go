package config

import (
	"github.com/gin-gonic/gin"
	"github.com/stc-ds-databricks-go/orm/query"
	"net/http"
)

func (a *AppConfig) Account(c *gin.Context) {
	dotNumber := c.Query("dotnumber")
	if dotNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dotnumber is required"})
		return
	}
	q := query.Use(a.GormDB)
	//account := model.SfAccount{}
	//stmt := fmt.Sprintf("SELECT * FROM sf_account WHERE DOT_Number__c = %s", dotNumber)
	expr := q.SfAccount.DOTNumberC.Eq(dotNumber)
	account, err := q.SfAccount.Where(expr).First()
	//tx := q.SfAccount.UnderlyingDB().Raw(stmt).First(&account)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"account": account})
}

// ========================

//		//cols := core.StringSlice{"Id", "Name", "DOT_Number__c", "CRM_Account_Number__c"}
//		//query := fmt.Sprintf(`SELECT %s FROM dev_silver.crm.sf_account WHERE Dot_Number__c=%s`, cols.Join(""), dotNumber)
//		rows, err := app.SqlDB.QueryContext(context.Background(), query)
//		if err != nil {
//			err = fmt.Errorf("error querying accounts: %s", err)
//			fmt.Println(err)
//			return
//		}
//		defer rows.Close()
//		account := model.SfAccount{}
//		for rows.Next() {
//			err = rows.Scan(&account.ID, &account.Name, &account.DOTNumberC, &account.CRMAccountNumberC)
//			if err != nil {
//				return
//			}
//			fmt.Println(account)
//		}
