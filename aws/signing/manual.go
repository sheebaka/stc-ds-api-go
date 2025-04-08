package signing

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	localConfig "github.com/stc-ds-databricks-go/config"
	"net/url"
	"os"
	"time"
)

func NewSignerConfig() (signerConfig *SignerConfig, err error) {
	_ = godotenv.Load()
	signerConfig = &SignerConfig{}
	if err = signerConfig.ReadConfig(); err != nil {
		return
	}
	err = signerConfig.RetrieveCredentials()
	currentTime := time.Now().UTC()
	signerConfig.amzdate = currentTime.Format("20060102T150405Z")
	signerConfig.datestamp = currentTime.Format("20060102")
	return
}

type SignerConfig struct {
	*viper.Viper
	aws.Config
	aws.Credentials
	amzdate   string
	datestamp string
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
	c.Config, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(c.GetString("api.aws_profile")))
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

func (c *SignerConfig) BuildRequestAuthHeaders(payloadHash, algorithm, credentialScope, signedHeaders, signature string) (headers core.Map[string]) {
	authorizationHeader := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", algorithm, c.AccessKeyID, credentialScope, signedHeaders, signature)
	headers = core.Map[string]{
		"Authorization":        authorizationHeader,
		"x-amz-date":           c.amzdate,
		"x-amz-content-sha256": payloadHash,
	}
	if c.Credentials.SessionToken != "" {
		headers.Put("x-amz-security-token", c.SessionToken)
	}
	return
}

func (c *SignerConfig) CalculateSignature(stringToSign string) (signatureString string) {
	signingKey := getSignatureKey(c.SecretAccessKey, c.datestamp, c.Region, c.GetString("api.service"))
	signatureSHA := hmac.New(sha256.New, signingKey)
	signatureSHA.Write([]byte(stringToSign))
	signatureString = hex.EncodeToString(signatureSHA.Sum(nil))
	return
}

func (c *SignerConfig) GetStringToSign(canonicalRequest string) (stringToSign, algorithm, credentialScope string) {
	algorithm = "AWS4-HMAC-SHA256"
	credentialScope = core.NewStringSlice(c.datestamp, c.Region, c.GetString("api.service"), "aws4_request").Join("/")
	stringToSign = core.NewStringSlice(algorithm, c.amzdate, credentialScope, hashAndEncode(canonicalRequest)).Join("\n")
	return
}

func (c *SignerConfig) GetCanonicalRequest() (canonicalRequest, payloadHash, signedHeaders string) {
	canonicalQueryString := ""
	headers := core.NewMap[string]().Put("host", c.GetString("api.api_gateway_invoke_dns"))
	canonicalHeadersMap := headers.Put("x-amz-date", c.amzdate)
	if c.SessionToken != "" {
		canonicalHeadersMap.Put("x-amz-security-token", c.SessionToken)
	}
	ss := core.NewStringSlice()
	for _, key := range canonicalHeadersMap.KeysSorted() {
		ss = ss.AppendPtr(fmt.Sprintf("%s:%s", key, canonicalHeadersMap.Get(key)))
	}
	canonicalHeaders := ss.Join("\n")
	signedHeaders = canonicalHeadersMap.Keys().Sorted().Join(";")
	endpointURl, err := url.Parse(c.GetString("api.vpc_endpoint_dns"))
	if err != nil {
		return
	}
	if endpointURl.Path == "" {
		endpointURl.Path = "/development/"
	}
	payloadHash = hashAndEncode()
	canonicalUri := endpointURl.Path
	canonicalRequestSS := core.NewStringSlice("GET", canonicalUri, canonicalQueryString, canonicalHeaders, signedHeaders, payloadHash)
	canonicalRequest = canonicalRequestSS.Join("\n")
	return
}
