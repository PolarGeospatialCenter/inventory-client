package client

import (
	"fmt"
	"net/url"

	"github.com/PolarGeospatialCenter/vaulthelper/pkg/vaulthelper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

type InventoryApi struct {
	AwsConfigs []*aws.Config
	BaseUrl    *url.URL
}

func NewInventoryApiDefaultConfig(profile string) (*InventoryApi, error) {
	if profile == "" {
		profile = "default"
	}

	cfg := viper.New()
	cfg.SetConfigName(profile)
	cfg.AddConfigPath("/etc/inventory")
	cfg.AddConfigPath("$HOME/.inventory")
	err := cfg.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading ~/.inventory/default.yml configuration file: %v", err)
	}

	baseUrl, err := url.Parse(cfg.GetString("baseurl"))
	if err != nil {
		return nil, fmt.Errorf("unable to parse base url: %v", err)
	}

	awsConfig := &aws.Config{}
	awsConfig.WithRegion(cfg.GetString("aws.region"))

	if vault_role := cfg.GetString("aws.vault_role"); vault_role != "" {
		vaultClient, err := vaulthelper.NewClient(vault.DefaultConfig())
		if err != nil {
			return nil, fmt.Errorf("unable to connect to vault: %v", err)
		}
		credProvider := &vaulthelper.VaultAwsStsCredentials{
			VaultClient: vaultClient,
			VaultRole:   vault_role,
		}
		awsConfig.WithCredentials(credentials.NewCredentials(credProvider))
	}

	if awsProfile := cfg.GetString("aws.profile"); awsProfile != "" {
		awsConfig.WithCredentials(credentials.NewSharedCredentials("", awsProfile))
	}

	return NewInventoryApi(baseUrl, awsConfig), nil

}

func NewInventoryApi(baseUrl *url.URL, configs ...*aws.Config) *InventoryApi {
	baseUrl, _ = baseUrl.Parse(baseUrl.Path + "/")
	return &InventoryApi{BaseUrl: baseUrl, AwsConfigs: configs}
}

func (i *InventoryApi) Url(endpointPath string) string {
	if endpointPath[0] == '/' {
		endpointPath = endpointPath[1:]
	}
	u, _ := i.BaseUrl.Parse(endpointPath)
	return u.String()
}

func (i *InventoryApi) Node() *Node {
	return &Node{Inventory: i}
}

func (i *InventoryApi) NodeConfig() *NodeConfig {
	return &NodeConfig{Inventory: i}
}

func (i *InventoryApi) Network() *Network {
	return &Network{Inventory: i}
}

func (i *InventoryApi) System() *System {
	return &System{Inventory: i}
}
