package client

import (
	"fmt"
	"net"
	"net/http"

	"github.com/PolarGeospatialCenter/inventory/pkg/inventory/types"
)

type IPAM struct {
	Inventory *InventoryApi
}

func (i *IPAM) GetIPReservation(ip net.IP) (*types.IPReservation, error) {
	client := NewRestClient(i.Inventory.AwsConfigs...)

	response, err := client.Client().NewRequest().Execute(http.MethodGet, i.Inventory.Url(fmt.Sprintf("/ipam/ip/%s", ip.String())))
	if err != nil {
		return nil, fmt.Errorf("unable to get ip reservation: %v", err)
	}
	reservation := &types.IPReservation{}
	err = UnmarshalApiResponse(response, reservation)
	return reservation, err
}

func (i *IPAM) GetIPReservationsByMAC(mac net.HardwareAddr) (*types.IPReservation, error) {
	client := NewRestClient(i.Inventory.AwsConfigs...)

	url := i.Inventory.Url(fmt.Sprintf("/ipam/ip?mac=%s", mac.String()))

	response, err := client.Client().NewRequest().Execute(http.MethodGet, url)
	if err != nil {
		return nil, fmt.Errorf("unable to get ip reservation: %v", err)
	}
	reservations := &types.IPReservationList{}
	err = UnmarshalApiResponse(response, reservations)
	return reservations, err
}

func (i *IPAM) UpdateIPReservation(modified *types.IPReservation) (*types.IPReservation, error) {
	client := NewRestClient(i.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()
	request.SetBody(modified)

	response, err := request.Execute(http.MethodPut, i.Inventory.Url(fmt.Sprintf("/ipam/ip/%s", modified.IP.IP.String())))
	if err != nil {
		return nil, fmt.Errorf("unable to get ip reservation: %v", err)
	}

	reservation := &types.IPReservation{}
	err = UnmarshalApiResponse(response, reservation)
	return reservation, err
}

func (i *IPAM) DeleteIPReservation(reservation *types.IPReservation) error {
	client := NewRestClient(i.Inventory.AwsConfigs...)

	response, err := client.Client().NewRequest().Execute(http.MethodDelete, i.Inventory.Url(fmt.Sprintf("/ipam/ip/%s", reservation.IP.IP.String())))
	if err != nil {
		return fmt.Errorf("unable to get ip reservation: %v", err)
	}

	err = UnmarshalApiResponse(response, nil)
	return err
}

func (i *IPAM) CreateIPReservation(new *types.IpamIpRequest, ip net.IP) (*types.IPReservation, error) {
	client := NewRestClient(i.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()
	request.SetBody(new)

	url := "/ipam/ip"
	if ip != nil {
		url = fmt.Sprintf("/ipam/ip/%s", ip.String())
	}
	response, err := request.Execute(http.MethodPost, i.Inventory.Url(url))
	if err != nil {
		return nil, fmt.Errorf("unable to create ip reservation: %v", err)
	}

	reservation := &types.IPReservation{}
	err = UnmarshalApiResponse(response, reservation)
	return reservation, err
}
