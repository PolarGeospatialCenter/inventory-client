package client

import (
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	gock "gopkg.in/h2non/gock.v1"
)

func TestIAMAuth(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()

	config := &aws.Config{}
	config.WithRegion("us-east-2")
	config.WithCredentials(credentials.NewStaticCredentials("asdf", "secret", "token"))

	client := NewRestClient(config)

	gock.New("https://api.local/v1").Get("/foo").HeaderPresent("X-Amz-Date").HeaderPresent("Authorization")

	_, err := client.Client().NewRequest().Execute(http.MethodGet, "https://api.local/v1/foo")
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}

}
