package rancher

import (
	"fmt"

	"github.com/rancher/go-rancher/client"
)

type Client struct {
	RancherClient *client.RancherClient
}

// NewClient grabs config necessary and sets an inited client or returns an error
func NewClient(cattleURL string, cattleAccessKey string, cattleSecretKey string) (*Client, error) {
	apiClient, err := client.NewRancherClient(&client.ClientOpts{
		Url:       cattleURL,
		AccessKey: cattleAccessKey,
		SecretKey: cattleSecretKey,
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		RancherClient: apiClient,
	}, nil
}

// TODO: Add test and validate that queried service is in state to receive finish upgrade
func (cli *Client) FinishServiceUpgrade(serviceName string) (*client.Service, error) {
	filters := make(map[string]interface{})
	filters["name"] = serviceName
	services, err := cli.RancherClient.Service.List(&client.ListOpts{
		Filters: filters,
	})

	if err != nil {
		return nil, err
	}

	if len(services.Data) != 1 {
		panic("more than 1 service found for same name")
	}
	service, err := cli.RancherClient.Service.ActionFinishupgrade(&services.Data[0])

	return service, err
}

func (cli *Client) ServiceByName(name string) (*client.Service, error) {
	filters := make(map[string]interface{})
	filters["name"] = name
	services, err := cli.RancherClient.Service.List(&client.ListOpts{
		Filters: filters,
	})

	if err != nil {
		return nil, err
	}

	if len(services.Data) != 1 {
		panic("more than 1 service found for same name")
	}
	return cli.RancherClient.Service.ById(services.Data[0].Id)
}

// TODO: Add integration test
// If updating just launch config image do not overwrite the secondary launch configs
func (cli *Client) UpgradeServiceVersion(serviceName, runtimeVersion string) error {
	service, err := cli.ServiceByName(serviceName)

	if err != nil {
		return err
	}

	service.LaunchConfig.ImageUuid = fmt.Sprintf("docker:%s", runtimeVersion)

	serviceUpgrade := &client.ServiceUpgrade{
		Resource: client.Resource{},
		InServiceStrategy: &client.InServiceUpgradeStrategy{
			LaunchConfig: service.LaunchConfig,
		},
	}
	_, err = cli.RancherClient.Service.ActionUpgrade(service, serviceUpgrade)

	return err
}

// TODO: Add integration test
func (cli *Client) UpgradeServiceCodeVersion(serviceName, codeVersion string) error {
	service, err := cli.ServiceByName(serviceName)

	if err != nil {
		return err
	}

	service.SecondaryLaunchConfigs[0].(map[string]interface{})["imageUuid"] = fmt.Sprintf("docker:%s", codeVersion)

	serviceUpgrade := &client.ServiceUpgrade{
		Resource: client.Resource{},
		InServiceStrategy: &client.InServiceUpgradeStrategy{
			SecondaryLaunchConfigs: service.SecondaryLaunchConfigs,
		},
	}
	_, err = cli.RancherClient.Service.ActionUpgrade(service, serviceUpgrade)

	return err
}
