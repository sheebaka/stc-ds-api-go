package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

// const Profile = "pipeline-kafka-deploy"

const Profile = "pipeline-kafka-deploy"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(Profile))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	url := "https://7ld4iwn31l.execute-api.us-east-2.amazonaws.com/development/api/v1/carrier_status/account?dot_number=3134772"
	// The signer requires a payload hash. This hash is for an empty payload.
	hash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	signer := v4.NewSigner()
	err = signer.SignHTTP(context.TODO(), creds, req, hash, "execute-api", cfg.Region, time.Now())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
}
