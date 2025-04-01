package config

import (
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
)

// ========================

func ConfigureApp(driverName string) (app *AppConfig, err error) {
	err = core.ReadYamlFile(JoinRoot("config.yaml"), &app)
	if err != nil {
		err = fmt.Errorf("failed to read config: %w", err)
		return
	}
	sourceConfig, ok := app.FindAndGet(driverName)
	if ok && sourceConfig != nil {
		sourceConfig.Driver = driverName
		if sourceConfig.Secret != "" {
			if err = sourceConfig.getSecret(); err != nil {
				fmt.Println(err)
			}
		}
		if err = sourceConfig.resolve(); err != nil {
			err = fmt.Errorf("failed to configure database dsn: %v", err)
			return
		}
	}
	app.SourceConfig = sourceConfig
	if err = app.ConfigGormDB(); err != nil {
		err = fmt.Errorf("failed to get *gorm.DB: %w", err)
		return
	}
	return
}
