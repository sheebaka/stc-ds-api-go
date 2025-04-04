package setup

import (
	"database/sql"
	"fmt"
	"github.com/stc-ds-databricks-go/config"
)

func Setupdb(driver string) (app *config.AppConfig) {
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
