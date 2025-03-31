package main

import (
	"fmt"
	"github.com/stc-ds-databricks-go/config"
	"github.com/stc-ds-databricks-go/orm"
)

func main() {
	app, err := config.ConfigureApp(config.Postgres)
	if err != nil {
		fmt.Println(err)
		return
	}
	account, err := orm.FilterByDotNumber(app.GormDb, "3134772")
	if err != nil {
		fmt.Println(err)
		return
	}
	config.PrettyPrint(account)
}
