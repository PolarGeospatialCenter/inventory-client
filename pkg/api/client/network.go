package client

import (
	"fmt"
	"net/http"

	"github.com/PolarGeospatialCenter/inventory/pkg/inventory/types"
)

type Network struct {
	Inventory *InventoryApi
}

func (n *Network) Get(id string) (*types.Network, error) {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	response, err := client.Client().NewRequest().Execute(http.MethodGet, n.Inventory.Url(fmt.Sprintf("/network/%s", id)))
	if err != nil {
		return nil, fmt.Errorf("unable to get network: %v", err)
	}
	network := &types.Network{}
	err = UnmarshalApiResponse(response, network)
	return network, err
}

func (n *Network) GetAll() ([]*types.Network, error) {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	response, err := client.Client().NewRequest().Execute(http.MethodGet, n.Inventory.Url("/network"))
	if err != nil {
		return nil, fmt.Errorf("unable to get networks: %v", err)
	}
	networks := []*types.Network{}
	err = UnmarshalApiResponse(response, &networks)
	return networks, err
}

func (n *Network) Create(network *types.Network) error {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()
	request.SetBody(network)

	response, err := request.Execute(http.MethodPost, n.Inventory.Url("/network"))
	if err != nil {
		return fmt.Errorf("unable to get networks: %v", err)
	}

	return UnmarshalApiResponse(response, nil)
}

func (n *Network) Update(network *types.Network) error {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()
	request.SetBody(network)

	response, err := request.Execute(http.MethodPut, n.Inventory.Url(fmt.Sprintf("/network/%s", network.ID())))
	if err != nil {
		return fmt.Errorf("unable to update networks: %v", err)
	}

	return UnmarshalApiResponse(response, nil)
}

func (n *Network) Delete(network *types.Network) error {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()

	response, err := request.Execute(http.MethodDelete, n.Inventory.Url(fmt.Sprintf("/network/%s", network.ID())))
	if err != nil {
		return fmt.Errorf("unable to update networks: %v", err)
	}

	return UnmarshalApiResponse(response, nil)
}
