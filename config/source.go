package config

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gen/helper"
	"gorm.io/gorm"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	//
	_ "github.com/databricks/databricks-sql-go"
	_ "github.com/lib/pq"
)

func (x *SourceConfig) DSN() string {
	return x.DataSourceName
}

func (x *SourceConfig) resolve() (err error) {
	var t *template.Template
	buf, err := yaml.Marshal(x)
	if err != nil {
		err = fmt.Errorf("failed to marshal config: %w", err)
		return
	}
	t, err = template.New("dsn").Parse(string(buf))
	if err != nil {
		return
	}
	buffer := bytes.NewBuffer(nil)
	if err = t.Execute(buffer, x); err != nil {
		err = fmt.Errorf("failed to execute template: %w", err)
		return
	}
	err = yaml.Unmarshal(buffer.Bytes(), &x)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal config: %w", err)
		return
	}
	return
}

func (x *SourceConfig) getSecret() (err error) {
	t, err := template.New("secret").Parse(x.Secret)
	if err != nil {
		return
	}
	buf := bytes.NewBufferString("")
	err = t.Execute(buf, x)
	if err != nil {
		return
	}
	x.Secret = buf.String()
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(x.Env))
	if err != nil {
		err = fmt.Errorf("failed to load config: %w", err)
		return
	}
	client := secretsmanager.NewFromConfig(cfg)
	secretValueOut, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{SecretId: aws.String(x.Secret)})
	if err != nil {
		if strings.Contains(err.Error(), "unable to refresh SSO token") {
			cmd := exec.Command("aws", strings.Split("sso login --sso-session aws-sso", " ")...)
			cmd.Stdout = os.Stdout
			if err = cmd.Run(); err != nil {
				return
			}
			err = x.getSecret()
		}
	}
	if secretValueOut != nil {
		fmt.Printf("secret: %s\n", *secretValueOut.SecretString)
		err = json.Unmarshal([]byte(*secretValueOut.SecretString), &x)
		host := regexp.MustCompile(`.*//(.*):([0-9]{4})/([a-z]+)`).ReplaceAllString(x.Host, "$1")
		x.Host = host
	}
	return
}

func (x *SourceConfig) ConfigGormDB() (err error) {
	dsn := x.DSN()
	db, err := sql.Open(x.Driver, dsn)
	if err != nil {
		err = fmt.Errorf("failed to connect database: %w", err)
		return
	}
	dialector := postgres.New(postgres.Config{
		DriverName: x.Driver,
		Conn:       db,
	})
	x.GormDB, err = gorm.Open(dialector, &gorm.Config{})
	return
}

func (x *SourceConfig) StructName() string {
	return x.DriverConfig.Tables[0]
}

func (x *SourceConfig) TableName() string {
	return x.DriverConfig.Tables[0]
}

func (x *SourceConfig) Table() string {
	return x.DriverConfig.Tables[0]
}

func (x *SourceConfig) FileName() string {
	return x.TableName()
}

func (x *SourceConfig) ImportPkgPaths() (ss []string) {
	return
}

func (x *SourceConfig) Fields() (fs []helper.Field) {
	return
}
