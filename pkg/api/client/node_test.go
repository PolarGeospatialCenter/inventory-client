package client

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/PolarGeospatialCenter/inventory/pkg/inventory/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	gock "gopkg.in/h2non/gock.v1"
)

func TestNodeGet(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Get("node/test-000").
		Reply(http.StatusOK).
		BodyString(`{"InventoryID": "test-000"}`)

	node, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Node().Get("test-000")
	if err != nil {
		t.Errorf("unable to get all nodes: %v", err)
	}

	if node.ID() != "test-000" {
		t.Errorf("got wrong inventory id: %s", node.ID())
	}
}

func TestNodeGetAll(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Get("node").
		Reply(http.StatusOK).
		BodyString(`[{"InventoryID": "test-000"}]`)

	nodes, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Node().GetAll()
	if err != nil {
		t.Errorf("unable to get all nodes: %v", err)
	}

	if len(nodes) != 1 {
		t.Errorf("got the wrong number of nodes: expected 1, got %d", len(nodes))
	}

	if nodes[0].InventoryID != "test-000" {
		t.Errorf("got wrong inventory id: %s", nodes[0].InventoryID)
	}
}

func TestNodeCreateSuccess(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Post("node").
		BodyString(`{"InventoryID":"test-001","ChassisSubIndex":"","Tags":null,"Networks":null,"Role":"","Environment":"","System":"","Metadata":null,"LastUpdated":"0001-01-01T00:00:00Z"}`).
		Reply(http.StatusCreated)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Node().Create(&types.Node{InventoryID: "test-001"})
	if err != nil {
		t.Errorf("unable to create node: %v", err)
	}
}

func TestNodeCreateConflict(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Post("node").
		BodyString(`{"InventoryID":"test-001","ChassisSubIndex":"","Tags":null,"Networks":null,"Role":"","Environment":"","System":"","Metadata":null,"LastUpdated":"0001-01-01T00:00:00Z"}`).
		Reply(http.StatusConflict)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Node().Create(&types.Node{InventoryID: "test-001"})
	if err == nil {
		t.Errorf("no error returned for conflict")
	}
}

func TestNodeUpdateSuccess(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	node := &types.Node{InventoryID: "test-001", Role: "test-role"}
	gock.New(testBaseUrl.String()).
		Put("node/test-001").
		BodyString(`{"InventoryID":"test-001","ChassisSubIndex":"","Tags":null,"Networks":null,"Role":"test-role","Environment":"","System":"","Metadata":null,"LastUpdated":"0001-01-01T00:00:00Z"}`).
		Reply(http.StatusOK)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Node().Update(node)
	if err != nil {
		t.Errorf("unable to update node: %v", err)
	}
}

func TestNodeDeleteSuccess(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	node := &types.Node{InventoryID: "test-001", Role: "test-role"}
	gock.New(testBaseUrl.String()).
		Delete("node/test-001").
		Reply(http.StatusOK)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Node().Delete(node)
	if err != nil {
		t.Errorf("unable to delete node: %v", err)
	}
}

func TestNodeDeleteNotFound(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	node := &types.Node{InventoryID: "test-001", Role: "test-role"}
	gock.New(testBaseUrl.String()).
		Delete("node/test-001").
		Reply(http.StatusNotFound).BodyString(`{"status": "Not Found", "error": "object not found"}`)

	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).Node().Delete(node)
	if err == nil {
		t.Errorf("no error returned when not found returned")
	}
}
