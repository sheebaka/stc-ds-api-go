package main

import (
	"database/sql"
	"fmt"
	_ "github.com/databricks/databricks-sql-go"
	"github.com/gin-gonic/gin"
	"github.com/stc-ds-databricks-go/config"
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
	app.SqlDB, err = sql.Open("databricks", dsn)
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
	carrier.GET("/account", app.Account)
	err = router.Run()
	return
}
