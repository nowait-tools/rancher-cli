package rancher

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nowait/rancher-cli/rancher/config"
	"github.com/rancher/go-rancher/client"
)

const (
	SERVICE_TYPE_SERVICE = "service"
)

var (
	upgradePollInterval = 150 * time.Millisecond
)

type Client struct {
	RancherClient *client.RancherClient
	Validator     config.Validator
}

type UpgradeOpts struct {
	Wait        bool
	ServiceLike string
	Service     string
	CodeTag     string
	RuntimeTag  string
	Interval    time.Duration
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

func (cli *Client) FinishServiceUpgrade(serviceName string) (*client.Service, error) {
	filters := make(map[string]interface{})
	filters["name"] = serviceName
	services, err := cli.RancherClient.Service.List(&client.ListOpts{
		Filters: filters,
	})

	if err != nil {
		return nil, err
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

	return &services.Data[0], nil
}

func (cli *Client) ServiceLikeName(likeName string) (services *client.ServiceCollection, err error) {
	filters := make(map[string]interface{})
	filters["name_like"] = getServiceLikeQuery(likeName)
	// Do not include service load balancers
	filters["kind"] = SERVICE_TYPE_SERVICE
	// TODO: Should filter include single environment.
	// If all users use environment specific keys that is fine
	// if they don't it could update multiple environments.
	services, err = cli.RancherClient.Service.List(&client.ListOpts{
		Filters: filters,
	})
	fmt.Printf("Upgrading %d services\n", len(services.Data))
	return
}

func (cli *Client) UpgradeServiceVersion(opts UpgradeOpts) error {
	service, err := cli.ServiceByName(opts.Service)

	if err != nil {
		return err
	}

	service.LaunchConfig.ImageUuid = fmt.Sprintf("docker:%s", opts.RuntimeTag)

	serviceUpgrade := &client.ServiceUpgrade{
		Resource: client.Resource{},
		InServiceStrategy: &client.InServiceUpgradeStrategy{
			// TODO: this should used to value passed in from the UpgradeOpts
			LaunchConfig:   service.LaunchConfig,
			BatchSize:      1,
			IntervalMillis: 10000,
			StartFirst:     true,
		},
	}
	_, err = cli.RancherClient.Service.ActionUpgrade(service, serviceUpgrade)

	return err
}

func (cli *Client) UpgradeServiceCodeVersion(opts UpgradeOpts) error {
	service, err := cli.ServiceByName(opts.Service)

	if err != nil {
		return err
	}

	serviceUpgrade := updateCodeImage(service, opts.CodeTag)
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
				err = Wait(cli, service, opts)
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
			if result.Error != nil {
				// Rollback upgrade, it failed
				failed = true
				fmt.Printf("service with name %s failed with: %v\n", result.Service.Name, result.Error)
				if opts.Wait {
					_, err := cli.RancherClient.Service.ActionCancelupgrade(result.Service)

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
			// TODO: this should used to value passed in from the UpgradeOpts
			BatchSize:              1,
			IntervalMillis:         10000,
			StartFirst:             true,
			SecondaryLaunchConfigs: service.SecondaryLaunchConfigs,
		},
	}
}

func Wait(cli *Client, srv *client.Service, opts UpgradeOpts) error {
	ch := make(chan error)
	go func() {
		<-time.After(opts.Interval * 20)
		ch <- errors.New("finishing upgrade timed out")
	}()
	go func() {
		for {
			if srv.Transitioning != "yes" {
				ch <- nil
			}

			time.Sleep(upgradePollInterval)

			err := cli.RancherClient.Reload(&srv.Resource, srv)
			if err != nil {
				ch <- err
			}
		}
	}()

	return <-ch
}

func getServiceLikeQuery(serviceName string) string {
	return serviceName + "%"
}
