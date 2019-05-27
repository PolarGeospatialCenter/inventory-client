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

type InventoryApiConfig struct {
	BaseURL string
	Aws     Aws
}

type Aws struct {
	Region    string
	VaultRole string `mapstructure:"vault_role"`
	Profile   string
}

// NewInventoryApiFromConfig returns a new inventory API from the config passed in.
func NewInventoryApiFromConfig(cfg *InventoryApiConfig) (*InventoryApi, error) {

	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse base url: %v", err)
	}

	awsConfig := &aws.Config{}
	awsConfig.WithRegion(cfg.Aws.Region)

	if vaultRole := cfg.Aws.VaultRole; vaultRole != "" {
		vaultClient, err := vaulthelper.NewClient(vault.DefaultConfig())
		if err != nil {
			return nil, fmt.Errorf("unable to connect to vault: %v", err)
		}
		credProvider := &vaulthelper.VaultAwsStsCredentials{
			VaultClient: vaultClient,
			VaultRole:   vaultRole,
		}
		awsConfig.WithCredentials(credentials.NewCredentials(credProvider))
	}

	if awsProfile := cfg.Aws.Profile; awsProfile != "" {
		awsConfig.WithCredentials(credentials.NewSharedCredentials("", awsProfile))
	}

	return NewInventoryApi(baseURL, awsConfig), nil

}

// NewInventoryApiDefaultConfig returns an InventoryApi based on config read by viper
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

	inventoryConfig := &InventoryApiConfig{}

	err = cfg.Unmarshal(inventoryConfig)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling InventoryApiConfig: %v", err)
	}

	return NewInventoryApiFromConfig(inventoryConfig)

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

func (i *InventoryApi) IPAM() *IPAM {
	return &IPAM{Inventory: i}
}
