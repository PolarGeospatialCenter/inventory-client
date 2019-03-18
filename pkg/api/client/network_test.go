package client

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/PolarGeospatialCenter/inventory/pkg/inventory/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	gock "gopkg.in/h2non/gock.v1"
)

func TestNetworkGet(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Get("network/test-000").
		Reply(http.StatusOK).
		BodyString(`{"Name": "test-000"}`)

	network, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Network().Get("test-000")
	if err != nil {
		t.Errorf("unable to get all networks: %v", err)
	}

	if network.ID() != "test-000" {
		t.Errorf("got wrong inventory id: %s", network.ID())
	}
}

func TestNetworkGetAll(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Get("network").
		Reply(http.StatusOK).
		BodyString(`[{"Name": "test-000"}]`)

	networks, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Network().GetAll()
	if err != nil {
		t.Errorf("unable to get all networks: %v", err)
	}

	if len(networks) != 1 {
		t.Errorf("got the wrong number of networks: expected 1, got %d", len(networks))
	}

	if networks[0].ID() != "test-000" {
		t.Errorf("got wrong inventory id: %s", networks[0].ID())
	}
}

func TestNetworkCreateSuccess(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	network := &types.Network{Name: "test-001"}
	expectedBody, err := json.Marshal(network)
	if err != nil {
		t.Errorf("unable to marshal updated network for testing: %v", err)
	}

	gock.New(testBaseUrl.String()).
		Post("network").BodyString(string(expectedBody)).
		Reply(http.StatusCreated)

	err = NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Network().Create(network)
	if err != nil {
		t.Errorf("unable to create network: %v", err)
	}
}

func TestNetworkCreateConflict(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Post("network").
		Reply(http.StatusConflict)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Network().Create(&types.Network{Name: "test-001"})
	if err == nil {
		t.Errorf("no error returned for conflict")
	}
}

func TestNetworkUpdateSuccess(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	network := &types.Network{Name: "test-001", Metadata: types.Metadata{"foo": "bar"}}
	expectedBody, err := json.Marshal(network)
	if err != nil {
		t.Errorf("unable to marshal updated network for testing: %v", err)
	}

	gock.New(testBaseUrl.String()).
		Put("network/test-001").BodyString(string(expectedBody)).
		Reply(http.StatusOK)

	err = NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Network().Update(network)
	if err != nil {
		t.Errorf("unable to update network: %v", err)
	}
}

func TestNetworkDeleteSuccess(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	network := &types.Network{Name: "test-001"}
	gock.New(testBaseUrl.String()).
		Delete("network/test-001").
		Reply(http.StatusOK)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Network().Delete(network)
	if err != nil {
		t.Errorf("unable to delete network: %v", err)
	}
}

func TestNetworkDeleteNotFound(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	network := &types.Network{Name: "test-001"}
	gock.New(testBaseUrl.String()).
		Delete("network/test-001").
		Reply(http.StatusNotFound).BodyString(`{"status": "Not Found", "error": "object not found"}`)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Network().Delete(network)
	if err == nil {
		t.Errorf("no error returned when not found returned")
	}
}
