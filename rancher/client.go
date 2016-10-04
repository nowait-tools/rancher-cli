package rancher

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/rancher/go-rancher/client"
)

type Client struct {
	RancherClient *client.RancherClient
}

type UpgradeOpts struct {
	Wait        bool
	ServiceLike string
	Service     string
	CodeTag     string
	RuntimeTag  string
}

type UpgradeResult struct {
	Service *client.Service
	Error   error
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

func (cli *Client) ServiceLikeName(likeName string) (services *client.ServiceCollection, err error) {
	filters := make(map[string]interface{})
	filters["name_like"] = likeName + "%"
	// Do not include service load balancers
	filters["kind"] = "service"
	// TODO: Might need to include environment id here.
	// If all users use environment specific keys that is fine
	// if they don't it might cause problems.
	services, err = cli.RancherClient.Service.List(&client.ListOpts{
		Filters: filters,
	})
	fmt.Printf("Upgrading %d services\n", len(services.Data))
	return
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
			// TODO: Figure out what the correct batch size and interval is.
			// Maybe this should be a configurable parameter?
			LaunchConfig:   service.LaunchConfig,
			BatchSize:      1,
			IntervalMillis: 10000,
			StartFirst:     true,
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

	serviceUpgrade := updateCodeImage(service, codeVersion)
	_, err = cli.RancherClient.Service.ActionUpgrade(service, serviceUpgrade)

	return err
}

func (cli *Client) UpgradeServiceWithNameLike(opts UpgradeOpts) error {
	failed := false
	services, err := cli.ServiceLikeName(opts.ServiceLike)

	if err != nil {
		return err
	}

	serviceCount := len(services.Data)
	upgradeErrs := make(chan UpgradeResult, serviceCount)

	for _, service := range services.Data {
		go func(srv client.Service) {
			serviceUpgrade := updateCodeImage(&srv, opts.CodeTag)
			service, err := cli.RancherClient.Service.ActionUpgrade(&srv, serviceUpgrade)

			if err != nil {
				upgradeErrs <- UpgradeResult{
					Service: &srv,
					Error:   err,
				}
				return
			}

			if opts.Wait {
				err = Wait(cli, service)
				if err == nil {
					_, err = cli.RancherClient.Service.ActionFinishupgrade(service)
				}
			}
			upgradeErrs <- UpgradeResult{
				Service: service,
				Error:   err,
			}
		}(service)
	}
	count := 0
	for {
		select {
		case result := <-upgradeErrs:
			log.Printf("%v", result)
			if result.Error != nil {
				// Rollback upgrade, it failed
				failed = true
				if opts.Wait {
					_, err := cli.RancherClient.Service.ActionRollback(result.Service)

					if err != nil {
						log.Fatalf("rollback failed with error: %v", err)
					}
				}
			}
			count++
			if count == serviceCount {
				if failed {
					return errors.New("upgrading services failed")
				}
				return nil
			}
		default:
			<-time.After(50 * time.Millisecond)
		}
	}
	return nil
}

func updateCodeImage(service *client.Service, codeVersion string) *client.ServiceUpgrade {
	service.SecondaryLaunchConfigs[0].(map[string]interface{})["imageUuid"] = fmt.Sprintf("docker:%s", codeVersion)

	return &client.ServiceUpgrade{
		Resource: client.Resource{},
		InServiceStrategy: &client.InServiceUpgradeStrategy{
			// TODO: Figure out what the correct batch size and interval is.
			// Maybe this should be a configurable parameter?
			BatchSize:              1,
			IntervalMillis:         10000,
			StartFirst:             true,
			SecondaryLaunchConfigs: service.SecondaryLaunchConfigs,
		},
	}
}

func Wait(cli *Client, srv *client.Service) error {
	for {
		if srv.Transitioning != "yes" {
			return nil
		}

		time.Sleep(150 * time.Millisecond)

		err := cli.RancherClient.Reload(&srv.Resource, srv)
		if err != nil {
			return err
		}
	}
}
