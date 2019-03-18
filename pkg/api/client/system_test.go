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

func TestSystemGet(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Get("system/test-000").
		Reply(http.StatusOK).
		BodyString(`{"Name": "test-000"}`)

	system, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).System().Get("test-000")
	if err != nil {
		t.Errorf("unable to get all systems: %v", err)
	}

	if system.ID() != "test-000" {
		t.Errorf("got wrong inventory id: %s", system.ID())
	}
}

func TestSystemGetAll(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Get("system").
		Reply(http.StatusOK).
		BodyString(`[{"Name": "test-000"}]`)

	systems, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).System().GetAll()
	if err != nil {
		t.Errorf("unable to get all systems: %v", err)
	}

	if len(systems) != 1 {
		t.Errorf("got the wrong number of systems: expected 1, got %d", len(systems))
	}

	if systems[0].ID() != "test-000" {
		t.Errorf("got wrong inventory id: %s", systems[0].ID())
	}
}

func TestSystemCreateSuccess(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	system := &types.System{Name: "test-001"}
	expectedBody, err := json.Marshal(system)
	if err != nil {
		t.Errorf("unable to marshal updated system for testing: %v", err)
	}

	gock.New(testBaseUrl.String()).
		Post("system").BodyString(string(expectedBody)).
		Reply(http.StatusCreated)

	err = NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).System().Create(system)
	if err != nil {
		t.Errorf("unable to create system: %v", err)
	}
}

func TestSystemCreateConflict(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Post("system").
		Reply(http.StatusConflict)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).System().Create(&types.System{Name: "test-001"})
	if err == nil {
		t.Errorf("no error returned for conflict")
	}
}

func TestSystemUpdateSuccess(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	system := &types.System{Name: "test-001", Metadata: types.Metadata{"foo": "bar"}}
	expectedBody, err := json.Marshal(system)
	if err != nil {
		t.Errorf("unable to marshal updated system for testing: %v", err)
	}

	gock.New(testBaseUrl.String()).
		Put("system/test-001").BodyString(string(expectedBody)).
		Reply(http.StatusOK)

	err = NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).System().Update(system)
	if err != nil {
		t.Errorf("unable to update system: %v", err)
	}
}

func TestSystemDeleteSuccess(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	system := &types.System{Name: "test-001"}
	gock.New(testBaseUrl.String()).
		Delete("system/test-001").
		Reply(http.StatusOK)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).System().Delete(system)
	if err != nil {
		t.Errorf("unable to delete system: %v", err)
	}
}

func TestSystemDeleteNotFound(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	system := &types.System{Name: "test-001"}
	gock.New(testBaseUrl.String()).
		Delete("system/test-001").
		Reply(http.StatusNotFound).BodyString(`{"status": "Not Found", "error": "object not found"}`)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).System().Delete(system)
	if err == nil {
		t.Errorf("no error returned when not found returned")
	}
}
