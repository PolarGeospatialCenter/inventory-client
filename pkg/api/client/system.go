package client

import (
	"fmt"
	"net/http"

	"github.com/PolarGeospatialCenter/inventory/pkg/inventory/types"
)

type System struct {
	Inventory *InventoryApi
}

func (n *System) Get(id string) (*types.System, error) {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	response, err := client.Client().NewRequest().Execute(http.MethodGet, n.Inventory.Url(fmt.Sprintf("/system/%s", id)))
	if err != nil {
		return nil, fmt.Errorf("unable to get system: %v", err)
	}
	system := &types.System{}
	err = UnmarshalApiResponse(response, system)
	return system, err
}

func (n *System) GetAll() ([]*types.System, error) {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	response, err := client.Client().NewRequest().Execute(http.MethodGet, n.Inventory.Url("/system"))
	if err != nil {
		return nil, fmt.Errorf("unable to get systems: %v", err)
	}
	systems := []*types.System{}
	err = UnmarshalApiResponse(response, &systems)
	return systems, err
}

func (n *System) Create(system *types.System) error {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()
	request.SetBody(system)

	response, err := request.Execute(http.MethodPost, n.Inventory.Url("/system"))
	if err != nil {
		return fmt.Errorf("unable to get systems: %v", err)
	}

	return UnmarshalApiResponse(response, nil)
}

func (n *System) Update(system *types.System) error {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()
	request.SetBody(system)

	response, err := request.Execute(http.MethodPut, n.Inventory.Url(fmt.Sprintf("/system/%s", system.ID())))
	if err != nil {
		return fmt.Errorf("unable to update systems: %v", err)
	}

	return UnmarshalApiResponse(response, nil)
}

func (n *System) Delete(system *types.System) error {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()

	response, err := request.Execute(http.MethodDelete, n.Inventory.Url(fmt.Sprintf("/system/%s", system.ID())))
	if err != nil {
		return fmt.Errorf("unable to update systems: %v", err)
	}

	return UnmarshalApiResponse(response, nil)
}
