package main

import (
	"database/sql"
	"fmt"
	_ "github.com/databricks/databricks-sql-go"
	"github.com/gin-gonic/gin"
	"github.com/stc-ds-databricks-go/config"
	"github.com/stc-ds-databricks-go/orm/model"
	"github.com/stc-ds-databricks-go/orm/query"
	"time"
)

const Databricks = "databricks"

func setupdb(driver string) (app *config.AppConfig) {
	var err error
	app, err = config.ConfigureApp(driver)
	if err != nil {
		err = fmt.Errorf("error configuring application: %s", err)
		fmt.Println(err)
		return
	}
	dsn := app.DSN()
	app.SqlDB, err = sql.Open("databricks", dsn)
	if err != nil {
		err = fmt.Errorf("error connecting to databricks: %s", err)
		fmt.Println(err)
	}
	return
}

func main() {
	app := setupdb(Databricks)
	err := Router(app)
	if err != nil {
		panic(err)
	}
}

func Router(app *config.AppConfig) (err error) {
	router := gin.Default()
	carrier := router.Group("/api/v1/carrier_status")
	carrier.GET("/account", func(c *gin.Context) {
		start := time.Now()
		dotNumber := c.Query("dotnumber")
		sfAccountDo := query.Use(app.GormDB).SfAccount
		var account *model.SfAccount
		account, err = sfAccountDo.Where(sfAccountDo.DOTNumberC.Eq(dotNumber)).First()

		stop := time.Now()
		fmt.Println(stop.Sub(start))
	})
}

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
