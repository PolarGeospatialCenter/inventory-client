package client

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/PolarGeospatialCenter/inventory/pkg/inventory/types"
)

type NodeConfig struct {
	Inventory *InventoryApi
}

func (c *NodeConfig) GetByMac(mac net.HardwareAddr) (*types.InventoryNode, error) {
	client := NewRestClient(c.Inventory.AwsConfigs...)
	response, err := client.Client().NewRequest().Execute(http.MethodGet, c.Inventory.Url(fmt.Sprintf("/nodeconfig?mac=%s", mac.String())))
	if err != nil {
		return nil, fmt.Errorf("unable to get node: %v", err)
	}

	node := []*types.InventoryNode{}
	err = UnmarshalApiResponse(response, &node)
	if err != nil {
		return nil, err
	}

	if len(node) != 1 {
		return nil, fmt.Errorf("unable to get node: unexpeceted number of nodes returned: %d", len(node))
	}
	return node[0], err
}

func (c *NodeConfig) Get(id string) (*types.InventoryNode, error) {
	client := NewRestClient(c.Inventory.AwsConfigs...)
	response, err := client.Client().NewRequest().Execute(http.MethodGet, c.Inventory.Url(fmt.Sprintf("/nodeconfig/%s", id)))
	if err != nil {
		return nil, fmt.Errorf("unable to get nodes: %v", err)
	}

	node := &types.InventoryNode{}
	err = UnmarshalApiResponse(response, node)
	return node, err
}

func (c *NodeConfig) GetAll() ([]*types.InventoryNode, error) {
	client := NewRestClient(c.Inventory.AwsConfigs...)
	response, err := client.Client().NewRequest().Execute(http.MethodGet, c.Inventory.Url("/nodeconfig"))
	if err != nil {
		return nil, fmt.Errorf("unable to get nodes: %v", err)
	}

	nodes := []*types.InventoryNode{}
	err = json.Unmarshal(response.Body(), &nodes)
	return nodes, err
}
