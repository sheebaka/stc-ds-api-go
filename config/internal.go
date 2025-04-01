package config

import (
	"encoding/json"
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// ========================

const (
	Databricks = "databricks"
	Postgres   = "postgres"
	DotNumber  = "3134772"
)

var Drivers = core.StringSlice{Databricks, Postgres}

var (
	_, b, _, _  = runtime.Caller(0)
	ProjectRoot = filepath.Join(filepath.Dir(b), "..")
)

func JoinRoot(elems ...string) string {
	return filepath.Join(ProjectRoot, filepath.Join(elems...))
}

// ========================

func ToSnakeCase(in string) (out string) {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(in, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	out = strings.ToLower(snake)
	return
}

func ToTitleCase(in string) (out string) {
	caser := cases.Title(language.English, cases.NoLower)
	for _, word := range strings.Split(in, "_") {
		word = caser.String(word)
		out += word
	}
	return
}

func Pprint(v interface{}) {
	x, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(x))
}
