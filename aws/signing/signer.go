package signing

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	localConfig "github.com/stc-ds-databricks-go/config"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type SignerConfig struct {
	*viper.Viper
	aws.Config
	aws.Credentials
	amzdate        string
	datestamp      string
	Headers        *core.Map[string]
	VpcEndpointUrl *url.URL
}

func NewSignerConfig() (signerConfig *SignerConfig, err error) {
	_ = godotenv.Load()
	signerConfig = &SignerConfig{}
	if err = signerConfig.ReadConfig(); err != nil {
		return
	}
	signerConfig.VpcEndpointUrl, err = url.Parse(signerConfig.GetString("vpc_endpoint_url"))
	err = signerConfig.RetrieveCredentials()
	currentTime := time.Now().UTC()
	signerConfig.amzdate = currentTime.Format("20060102T150405Z")
	signerConfig.datestamp = currentTime.Format("20060102")
	return
}

func (c *SignerConfig) ReadConfig() (err error) {
	viper.SetConfigType("yaml")
	fp := localConfig.JoinRoot("aws", "config.yaml")
	viper.SetConfigFile(fp)
	err = viper.ReadInConfig()
	c.Viper = viper.GetViper()
	return
}

func (c *SignerConfig) RetrieveCredentials() (err error) {
	c.Config, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(c.GetString("aws_profile")))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.Credentials, err = c.Config.Credentials.Retrieve(context.TODO())
	if s := os.Getenv("AWS_ACCESS_KEY_ID"); s != "" {
		c.Credentials.AccessKeyID = s
	}
	if s := os.Getenv("AWS_SECRET_ACCESS_KEY"); s != "" {
		c.Credentials.SecretAccessKey = s
	}
	if s := os.Getenv("AWS_SESSION_TOKEN"); s != "" {
		c.Credentials.SessionToken = s
	}
	return
}

func hashAndEncode(ss ...string) (out string) {
	var s string
	if len(ss) == 0 {
		s = ""
	}
	if len(ss) == 1 {
		s = ss[0]
	}
	h := sha256.New()
	h.Write([]byte(s))
	out = hex.EncodeToString(h.Sum(nil))
	return
}

func sign(key []byte, value string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(value))
	return h.Sum(nil)
}

func getSignatureKey(secretKey, credentialDate, region, service string) (signingKey []byte) {
	kSecret := []byte(fmt.Sprintf("AWS4%s", secretKey))
	kDate := sign(kSecret, credentialDate)
	kRegion := sign(kDate, region)
	kService := sign(kRegion, service)
	signingKey = sign(kService, "aws4_request")
	return
}

// ========================

// GetCanonicalRequest 1
func (c *SignerConfig) GetCanonicalRequest(headers core.Map[string]) (canonicalRequest, payloadHash, signedHeaders string) {
	canonicalQueryString := c.VpcEndpointUrl.RawQuery
	canonicalHeadersMap := headers.Clone()
	canonicalHeadersMap.Put("x-amz-date", c.amzdate)
	if c.SessionToken != "" {
		canonicalHeadersMap.Put("x-amz-security-token", c.SessionToken)
	}
	ss := core.NewStringSlice()
	for _, key := range canonicalHeadersMap.KeysSorted() {
		ss = ss.AppendPtr(fmt.Sprintf("%s:%s\n", key, canonicalHeadersMap.Get(key)))
	}
	canonicalHeaders := ss.Join("")
	signedHeaders = canonicalHeadersMap.Keys().Sorted().Join(";")
	payloadHash = hashAndEncode()
	canonicalUri := c.VpcEndpointUrl.Path
	canonicalRequestSS := core.NewStringSlice("GET", canonicalUri, canonicalQueryString, canonicalHeaders, signedHeaders, payloadHash)
	canonicalRequest = canonicalRequestSS.Join("\n")
	return
}

// GetStringToSign 2
func (c *SignerConfig) GetStringToSign(canonicalRequest string) (stringToSign, algorithm, credentialScope string) {
	algorithm = c.GetString("algorithm")
	credentialScope = core.NewStringSlice(c.datestamp, c.Region, c.GetString("service"), "aws4_request").Join("/")
	stringToSign = core.NewStringSlice(algorithm, c.amzdate, credentialScope, hashAndEncode(canonicalRequest)).Join("\n")
	return
}

// CalculateSignature 3
func (c *SignerConfig) CalculateSignature(stringToSign string) (signatureString string) {
	signingKey := getSignatureKey(c.SecretAccessKey, c.datestamp, c.Region, c.GetString("service"))
	signatureSHA := hmac.New(sha256.New, signingKey)
	signatureSHA.Write([]byte(stringToSign))
	signatureString = hex.EncodeToString(signatureSHA.Sum(nil))
	return
}

// BuildRequestAuthHeaders 4
func (c *SignerConfig) BuildRequestAuthHeaders(payloadHash, algorithm, credentialScope, signedHeaders, signature string) (headers *core.Map[string]) {
	authorizationHeader := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", algorithm, c.AccessKeyID, credentialScope, signedHeaders, signature)
	headers = core.NewMap[string]()
	headers.Put("Authorization", authorizationHeader)
	headers.Put("x-amz-date", c.amzdate)
	headers.Put("x-amz-content-sha256", payloadHash)
	//
	if c.Credentials.SessionToken != "" {
		headers.Put("x-amz-security-token", c.SessionToken)
	}
	return
}

// DoRequest 5
func (c *SignerConfig) DoRequest() (value core.Map[any], err error) {
	uri := c.VpcEndpointUrl.String()
	fmt.Printf("Creating new request with: %s\n", uri)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Host = c.GetString("invoke_url")
	fmt.Println(req.Host)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(buf, &value)
	return
}
