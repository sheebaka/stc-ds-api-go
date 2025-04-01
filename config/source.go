package config

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gen/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"text/template"
	//
	_ "github.com/databricks/databricks-sql-go"
	_ "github.com/lib/pq"
)

func (x *SourceConfig) DSN() string {
	return x.DataSourceName
}

func (x *SourceConfig) handleEnv() (err error) {
	_ = godotenv.Load()
	sm := core.Map[string]{}
	pre := strings.ToUpper(x.Driver)
	envVars := core.NewStringSlice(os.Environ()...).Sorted()
	for _, env := range envVars {
		ss := strings.Split(env, "=")
		if strings.HasPrefix(env, pre) {
			sm.Put(ss[0], ss[1])
			continue
		}
		others := Drivers.RemoveItems(x.Driver)
		unset := slices.ContainsFunc(others, func(s string) bool {
			return strings.HasPrefix(env, strings.ToUpper(s))
		})
		if !unset {
			continue
		}
		if err = os.Unsetenv(ss[0]); err != nil {
			err = fmt.Errorf("unable to unset env %s", ss[0])
			fmt.Println(err)
		}
		if os.Getenv(ss[0]) == "" {
			fmt.Printf("Unset env %s\n", ss[0])
		}
	}
	if !sm.IsEmpty() {
		fmt.Printf("found the following environment variables for %s driver:\n", x.Driver)
	}
	for k, v := range sm.Iter() {
		fmt.Printf("  ‣ %s = %s\n", k, v)
		if strings.Contains(k, "HOST") {
			_ = os.Setenv("HOST", v)
		}
	}
	if err = cleanenv.ReadEnv(x); err != nil {
		err = fmt.Errorf("unable to read env: %s", err)
	}
	return
}

func (x *SourceConfig) resolve() (err error) {
	if err = x.handleEnv(); err != nil {
		return
	}
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
	if x.Driver == Postgres {
		if x.SqlDB, err = sql.Open(x.Driver, dsn); err != nil {
			return
		}
		x.Dialector = postgres.New(postgres.Config{
			DriverName: x.Driver,
			Conn:       x.SqlDB,
			DSN:        dsn,
		})
	}
	if x.Driver == Databricks {
		if x.SqlDB, err = sql.Open(x.Driver, dsn); err != nil {
			return
		}
		x.Dialector = mysql.New(mysql.Config{
			DriverName: x.Driver,
			Conn:       x.SqlDB,
			DSN:        dsn,
		})
	}
	x.GormDB, err = gorm.Open(x.Dialector, &gorm.Config{
		NamingStrategy: x,
	})
	return
}

func (x *SourceConfig) JoinTableName(joinTable string) (s string) {
	return
}
func (x *SourceConfig) RelationshipFKName(schema.Relationship) (s string) { return }
func (x *SourceConfig) CheckerName(table string, column string) (s string) {
	return
}
func (x *SourceConfig) ColumnName(table string, column string) (s string) {
	return
}
func (x *SourceConfig) IndexName(table string, column string) (s string) {
	return
}
func (x *SourceConfig) UniqueName(table string, column string) (s string) {
	return
}
func (x *SourceConfig) SchemaDotTable() (s string) {
	return fmt.Sprintf("%s.%s", x.Schema, x.Table())
}
func (x *SourceConfig) StructName() string {
	return x.DriverConfig.Tables[0]
}
func (x *SourceConfig) TableName(table string) string {
	return x.DriverConfig.Tables[0]
}
func (x *SourceConfig) ModelName() string {
	return ToTitleCase(x.TableName(""))
}
func (x *SourceConfig) Table() string {
	return x.DriverConfig.Tables[0]
}
func (x *SourceConfig) SchemaName(schema string) string {
	return schema
}
func (x *SourceConfig) FileName() string {
	return x.TableName("")
}
func (x *SourceConfig) ImportPkgPaths() (ss []string) {
	return
}
func (x *SourceConfig) Fields() (fs []helper.Field) {
	return
}
