package main

import (
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"github.com/stc-ds-databricks-go/aws/signing"
	"net/url"
)

func main() {
	config, err := signing.NewSignerConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	invokeUrl, err := url.Parse(config.GetString("invoke_url"))
	if err != nil {
		fmt.Println(err)
		return
	}
	host := invokeUrl.Host
	if host == "" {
		host = invokeUrl.Path
	}
	headers := core.NewMap[string]().Put("host", host)
	// CanonicalRequest
	canonicalRequest, payloadHash, signedHeaders := config.GetCanonicalRequest(*headers)
	// GetStringToSign
	stringToSign, algorithm, credentialScope := config.GetStringToSign(canonicalRequest)
	// CalculateSignature
	signature := config.CalculateSignature(stringToSign)
	// BuildRequestAuthHeaders
	authHeaders := config.BuildRequestAuthHeaders(payloadHash, algorithm, credentialScope, signedHeaders, signature)
	// Add original Host header back
	config.Headers = headers.Add(*authHeaders)
	// DoRequest
	response, err := config.DoRequest()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Response:")
	fmt.Println(core.PrettyStruct(response))
}
