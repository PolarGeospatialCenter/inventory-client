package client

import (
	"fmt"
	"net/http"

	"github.com/PolarGeospatialCenter/inventory/pkg/inventory/types"
)

type Node struct {
	Inventory *InventoryApi
}

func (n *Node) Get(id string) (*types.Node, error) {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	response, err := client.Client().NewRequest().Execute(http.MethodGet, n.Inventory.Url(fmt.Sprintf("/node/%s", id)))
	if err != nil {
		return nil, fmt.Errorf("unable to get node: %v", err)
	}
	node := &types.Node{}
	err = UnmarshalApiResponse(response, node)
	return node, err
}

func (n *Node) GetAll() ([]*types.Node, error) {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	response, err := client.Client().NewRequest().Execute(http.MethodGet, n.Inventory.Url("/node"))
	if err != nil {
		return nil, fmt.Errorf("unable to get nodes: %v", err)
	}
	nodes := []*types.Node{}
	err = UnmarshalApiResponse(response, &nodes)
	return nodes, err
}

func (n *Node) Create(node *types.Node) error {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()
	request.SetBody(node)

	response, err := request.Execute(http.MethodPost, n.Inventory.Url("/node"))
	if err != nil {
		return fmt.Errorf("unable to get nodes: %v", err)
	}

	return UnmarshalApiResponse(response, nil)
}

func (n *Node) Update(node *types.Node) error {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()
	request.SetBody(node)

	response, err := request.Execute(http.MethodPut, n.Inventory.Url(fmt.Sprintf("/node/%s", node.ID())))
	if err != nil {
		return fmt.Errorf("unable to update nodes: %v", err)
	}

	return UnmarshalApiResponse(response, nil)
}

func (n *Node) Delete(node *types.Node) error {
	client := NewRestClient(n.Inventory.AwsConfigs...)

	request := client.Client().NewRequest()

	response, err := request.Execute(http.MethodDelete, n.Inventory.Url(fmt.Sprintf("/node/%s", node.ID())))
	if err != nil {
		return fmt.Errorf("unable to update nodes: %v", err)
	}

	return UnmarshalApiResponse(response, nil)
}
