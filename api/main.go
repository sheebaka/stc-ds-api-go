package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stc-ds-databricks-go/api/setup"
	"github.com/stc-ds-databricks-go/config"
	"github.com/stc-ds-databricks-go/logging"
	"github.com/stc-ds-databricks-go/orm/common"
	"net/http"
	"strconv"
)

var logger = logging.Logger()

func main() {
	app := setup.Setupdb(config.Databricks)
	err := Router(app)
	if err != nil {
		logger.Fatal(err)
	}
}

func Router(app *config.AppConfig) (err error) {
	router := gin.Default()
	if err = router.SetTrustedProxies(nil); err != nil {
		logger.Fatal(err)
		return
	}
	carrier := router.Group("/api/v1/carrier_status")
	carrier.GET("/account", func(c *gin.Context) {
		Account(app, c)
	})
	err = router.Run(app.GetPort())
	return
}

func Account(a *config.AppConfig, c *gin.Context) {
	logger.Info(c.FullPath())
	dotNumber := c.Query("dot_number")
	dotNumInt, err := strconv.Atoi(dotNumber)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("invalid dot number: %s", dotNumber)})
		return
	}
	result, err := common.FilterByDotNumber(a.GormDB, dotNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if result != nil && result.Active__c == "true" {
		result.Active__c = "active"
	}
	if result == nil {
		result = &common.Result{
			DOT_Number__c: int64(dotNumInt),
		}
	}
	c.JSON(http.StatusOK, result)
}
