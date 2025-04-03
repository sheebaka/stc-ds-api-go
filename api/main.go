package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stc-ds-databricks-go/config"
	"github.com/stc-ds-databricks-go/orm/common"
	"net/http"
	"strconv"
)

func setupdb(driver string) (app *config.AppConfig) {
	var err error
	app, err = config.ConfigureApp(driver)
	if err != nil {
		err = fmt.Errorf("error configuring application: %s", err)
		fmt.Println(err)
		return
	}
	dsn := app.DSN()
	app.SqlDB, err = sql.Open(config.Databricks, dsn)
	if err != nil {
		err = fmt.Errorf("error connecting to databricks: %s", err)
		fmt.Println(err)
	}
	return
}

func main() {
	app := setupdb(config.Databricks)
	err := Router(app)
	if err != nil {
		panic(err)
	}
}

func Router(app *config.AppConfig) (err error) {
	router := gin.Default()
	carrier := router.Group("/api/v1/carrier_status")
	carrier.GET("/account", func(c *gin.Context) {
		Account(app, c)
	})
	err = router.Run(app.GetPort())
	return
}

func Account(a *config.AppConfig, c *gin.Context) {
	fmt.Println(c.FullPath())
	dotNumber := c.Query("dotnumber")
	result, err := common.FilterByDotNumber(a.GormDB, dotNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//var result *common.Result
	//if len(results) == 1 {
	//	result = results[0]
	//	if result.Active__c == "true" {
	//		result.Active__c = "active"
	//	}
	//}
	if result != nil && result.Active__c == "true" {
		result.Active__c = "active"
	}
	dotNumInt, err := strconv.ParseInt(dotNumber, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if result == nil {
		result = &common.Result{
			DOT_Number__c: dotNumInt,
		}
	}
	c.JSON(http.StatusOK, result)
}
