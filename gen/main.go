package main

import (
	"fmt"
	"github.com/stc-ds-databricks-go/config"
)

const DriverName = config.Databricks

func main() {
	app, err := config.ConfigureApp(DriverName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = app.GenerateModel(); err != nil {
		fmt.Println(err)
		return
	}
}
