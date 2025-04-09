package main

import (
	"fmt"
	"github.com/stc-ds-databricks-go/aws/signing"
)

func main() {
	config, err := signing.NewSignerConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("AccessKeyId: %s\n", config.AccessKeyID)
	fmt.Printf("SecretAccessKey: %s\n", config.SecretAccessKey)
	fmt.Printf("SessionToken: %s\n", config.SessionToken)
	fmt.Println("========================")
	// CanonicalRequest
	canonicalRequest, payloadHash, signedHeaders := config.GetCanonicalRequest()
	fmt.Printf("canonicalRequest: %s\n", canonicalRequest)
	fmt.Println("========================")
	fmt.Printf("payloadHash: %s\n", payloadHash)
	fmt.Println("========================")
	fmt.Printf("signedHeaders: %s\n", signedHeaders)
	fmt.Println("========================")
	// GetStringToSign
	stringToSign, algorithm, credentialScope := config.GetStringToSign(canonicalRequest)
	fmt.Printf("stringToSign: %s\n", stringToSign)
	fmt.Println("========================")
	fmt.Printf("algorithm: %s\n", algorithm)
	fmt.Println("========================")
	fmt.Printf("credentialScope: %s\n", credentialScope)
	fmt.Println("========================")
	// CalculateSignature
	signature := config.CalculateSignature(stringToSign)
	fmt.Printf("signature: %s\n", signature)
	fmt.Println("========================")
	// BuildRequestAuthHeaders
	authHeaders := config.BuildRequestAuthHeaders(payloadHash, algorithm, credentialScope, signedHeaders, signature)
	fmt.Printf("authHeaders: %v\n", authHeaders)
	fmt.Println("========================")
	// Add original Host header back
	authHeaders.Put("Host", config.GetString("api.api_gateway_invoke_dns"))
	fmt.Printf("authHeaders: %v\n", authHeaders)
	fmt.Println("========================")
	// DoRequest
	config.Headers = authHeaders
	err = config.DoRequest()
	if err != nil {
		fmt.Println(err)
	}
}
