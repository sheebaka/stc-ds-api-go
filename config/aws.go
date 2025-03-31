package config

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AwsRootConfig struct {
	*http.Request
	accessKeyId     string `yaml:"access_key_id" env:"AWS_ACCESS_KEY_ID"`
	secretAccessKey string `yaml:"secret_access_key" env:"AWS_SECRET_ACCESS_KEY"`
	region          string `yaml:"region" env:"AWS_REGION"`
}

func NewAwsConfig() *AwsRootConfig {
	return &AwsRootConfig{}
}

func (ac *AwsRootConfig) buildRequest() {
	data := make(url.Values)
	data.Add("Action", "?")
	// add data for appropriate service ...
	// ...
	data.Add("AWSAccessKeyId", ac.accessKeyId)
	encodedData := data.Encode()
	body := strings.NewReader(encodedData)
	fmt.Println(body)
}

func (ac *AwsRootConfig) canonicalRequest(host string) (s string) {
	now := time.Now().UTC()
	date := now.Format("20060102T150405Z")
	//
	req := &http.Request{}
	ac.Request = req
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("host", host)
	req.Header.Set("x-amz-date", date)
	//
	//canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-date:%s\n", host, date)
	//signedHeaders := "host;x-amz-date"
	//
	//payloadHash := hashAndEncode()
	//
	return
}
