package main

import (
	"fmt"
	"github.com/ShamrockTrading/stc-ds-dataeng-go/core"
	"github.com/stc-ds-databricks-go/aws/signing"
	"io"
	"net/http"
)

func main() {
	config, err := signing.NewSignerConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	//config.Sign()
	headers := core.Map[string]{"Host": config.GetString("api.api_gateway_invoke_dns")}
	canonicalRequest, payloadHash, signedHeaders := config.GetCanonicalRequest()
	stringToSign, algorithm, credentialScope := config.GetStringToSign(canonicalRequest)
	signature := config.CalculateSignature(stringToSign)
	authHeaders := config.BuildRequestAuthHeaders(payloadHash, algorithm, credentialScope, signedHeaders, signature)
	headers.Add(authHeaders)
	fmt.Println(headers)
	uri := config.GetString("api.vpc_endpoint_dns")
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	fmt.Println(core.PrettyStruct(headers))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := io.ReadAll(res.Body)
	fmt.Println(string(b))
}
