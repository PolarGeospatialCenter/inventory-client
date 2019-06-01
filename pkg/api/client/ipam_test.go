package client

import (
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/PolarGeospatialCenter/inventory/pkg/inventory/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	gock "gopkg.in/h2non/gock.v1"
)

func TestIPReservationGet(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Get("ipam/ip/10.0.0.1").
		Reply(http.StatusOK).
		BodyString(`{"mac": "00:01:02:03:04:05","ip": "10.0.0.1/24","start": null,"end": null}`)

	reservation, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).IPAM().GetIPReservation(net.ParseIP("10.0.0.1"))
	if err != nil {
		t.Errorf("unable to get ip reservation: %v", err)
	}

	if reservation.MAC.String() != "00:01:02:03:04:05" {
		t.Errorf("mac address doesn't match expected value")
	}
	if reservation.Start != nil {
		t.Errorf("expected nil start time")
	}
}

func TestIPReservationGetByMAC(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Get("ipam/ip").
		MatchParam("mac", "00:01:02:03:04:05").
		Reply(http.StatusOK).
		BodyString(`[{"mac": "00:01:02:03:04:05","ip": "10.0.0.1/24","start": null,"end": null}]`)

	mac, _ := net.ParseMAC("00:01:02:03:04:05")
	reservation, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).
		IPAM().GetIPReservationsByMAC(mac)
	if err != nil {
		t.Errorf("unable to get ip reservation: %v", err)
	}

	if reservation[0].MAC.String() != "00:01:02:03:04:05" {
		t.Errorf("mac address doesn't match expected value")
	}
	if reservation[0].Start != nil {
		t.Errorf("expected nil start time")
	}
}

func TestIPReservationCreateSpecifiedIP(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Post("ipam/ip/10.0.0.2").
		Reply(http.StatusCreated).
		BodyString(`{
			"HostInformation": "",
			"start": "2019-05-27T02:22:30.97489765Z",
			"end": null,
			"ip": "10.0.0.2/24",
			"mac": "00:01:02:03:04:06",
			"gateway": "10.0.0.1",
			"dns": [
				"8.8.8.8",
				"1.1.1.1"
			]
		}`)

	mac, _ := net.ParseMAC("00:01:02:03:04:06")

	new := &types.IpamIpRequest{
		HwAddress: mac.String(),
		Subnet:    net.ParseIP("10.0.0.0").String(),
	}
	reservation, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).IPAM().CreateIPReservation(new, net.ParseIP("10.0.0.2"))
	if err != nil {
		t.Fatalf("unable to get ip reservation: %v", err)
	}

	if reservation.MAC.String() != "00:01:02:03:04:06" {
		t.Errorf("mac address doesn't match expected value")
	}
	if reservation.Start == nil {
		t.Errorf("expected non-nil start time")
	}
}

func TestIPReservationCreate(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Post("ipam/ip").
		Reply(http.StatusCreated).
		BodyString(`{
			"HostInformation": "",
			"start": "2019-05-27T02:22:30.97489765Z",
			"end": null,
			"ip": "10.0.0.2/24",
			"mac": "00:01:02:03:04:06",
			"gateway": "10.0.0.1",
			"dns": [
				"8.8.8.8",
				"1.1.1.1"
			]
		}`)

	mac, _ := net.ParseMAC("00:01:02:03:04:06")

	new := &types.IpamIpRequest{
		HwAddress: mac.String(),
		Subnet:    net.ParseIP("10.0.0.0").String(),
	}
	reservation, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).IPAM().CreateIPReservation(new, nil)
	if err != nil {
		t.Fatalf("unable to get ip reservation: %v", err)
	}

	if reservation.MAC.String() != "00:01:02:03:04:06" {
		t.Errorf("mac address doesn't match expected value")
	}
	if reservation.Start == nil {
		t.Errorf("expected non-nil start time")
	}
}
func TestIPReservationUpdate(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Put("ipam/ip/10.0.0.1").
		Reply(http.StatusOK).
		BodyString(`{"mac": "00:01:02:03:04:06","ip": "10.0.0.1/24","start": null,"end": null}`)

	mac, _ := net.ParseMAC("00:01:02:03:04:06")
	modified := &types.IPReservation{
		IP:  &net.IPNet{IP: net.ParseIP("10.0.0.1"), Mask: net.IPv4Mask(0xff, 0xff, 0xff, 0x00)},
		MAC: mac,
	}
	reservation, err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).IPAM().UpdateIPReservation(modified)
	if err != nil {
		t.Fatalf("unable to get ip reservation: %v", err)
	}

	if reservation.MAC.String() != "00:01:02:03:04:06" {
		t.Errorf("mac address doesn't match expected value")
	}
	if reservation.Start != nil {
		t.Errorf("expected nil start time")
	}
}
func TestIPReservationDelete(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Delete("ipam/ip/10.0.0.1").
		Reply(http.StatusOK)

	reservation := &types.IPReservation{IP: &net.IPNet{IP: net.ParseIP("10.0.0.1"), Mask: net.IPv4Mask(0xff, 0xff, 0xff, 0x00)}}
	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).IPAM().DeleteIPReservation(reservation)
	if err != nil {
		t.Errorf("unable to delete ip reservation: %v", err)
	}
}
func TestIPReservationDeleteNotFound(t *testing.T) {
	gock.DisableNetworking()
	defer gock.EnableNetworking()
	defer gock.Off()
	testBaseUrl, _ := url.Parse("https://inventory.api.local/v0/")

	gock.New(testBaseUrl.String()).
		Delete("ipam/ip/10.0.0.1").
		Reply(http.StatusNotFound).
		BodyString(`{"status": "Not Found", "error": "object not found"}`)

	reservation := &types.IPReservation{IP: &net.IPNet{IP: net.ParseIP("10.0.0.1"), Mask: net.IPv4Mask(0xff, 0xff, 0xff, 0x00)}}
	err := NewInventoryApi(testBaseUrl, &aws.Config{Credentials: credentials.NewStaticCredentials("id", "secret", "token")}).IPAM().DeleteIPReservation(reservation)
	if err == nil {
		t.Errorf("no error reported when deleting non-existent record")
	}
}
