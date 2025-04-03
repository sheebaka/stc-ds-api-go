package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stc-ds-databricks-go/config"
	"github.com/stc-ds-databricks-go/orm/common"
	"net/http"
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
	err = router.Run(":5555")
	return
}

func Account(a *config.AppConfig, c *gin.Context) {
	fmt.Println(c.FullPath())
	dotNumber := c.Query("dotnumber")
	results, err := common.FilterByDotNumber(a.GormDB, dotNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var result *common.Result
	if len(results) == 1 {
		result = results[0]
		if result.Active__c == "true" {
			result.Active__c = "active"
		}
	}
	c.JSON(http.StatusOK, result)
}
