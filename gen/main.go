package main

import (
	"fmt"
	"github.com/stc-ds-databricks-go/config"
)

func main() {
	app, err := config.ConfigureApp(config.Postgres)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = app.GenerateModel(); err != nil {
		fmt.Println(err)
		return
	}
}
