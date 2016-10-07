package rancher

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rancher/go-rancher/client"
)

var serviceName = "name"
var codeTag = "image-name:1.0"

type NoopService struct{}

func (srv *NoopService) List(opts *client.ListOpts) (*client.ServiceCollection, error) {
	return nil, nil
}
func (srv *NoopService) Create(opts *client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) Update(existing *client.Service, updates interface{}) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ById(id string) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) Delete(container *client.Service) error {
	return nil
}
func (srv *NoopService) ActionActivate(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionAddservicelink(*client.Service, *client.AddRemoveServiceLinkInput) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionCancelrollback(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionCancelupgrade(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionCreate(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionDeactivate(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionFinishupgrade(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionRemove(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionRemoveservicelink(*client.Service, *client.AddRemoveServiceLinkInput) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionRestart(*client.Service, *client.ServiceRestart) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionRollback(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *NoopService) ActionSetservicelinks(*client.Service, *client.SetServiceLinksInput) (*client.Service, error) {
	return nil, nil
}

func (srv *NoopService) ActionUpdate(*client.Service) (*client.Service, error) {
	return nil, nil
}

func (srv *NoopService) ActionUpgrade(*client.Service, *client.ServiceUpgrade) (*client.Service, error) {
	return nil, nil
}

type FailedService struct {
	ListOpts *client.ListOpts
	NoopService
}

func (srv *FailedService) List(opts *client.ListOpts) (*client.ServiceCollection, error) {
	return nil, errors.New("retrieving services failed with error")
}

type SuccessfulService struct {
	ListOpts *client.ListOpts
	NoopService
}

func (srv *SuccessfulService) List(opts *client.ListOpts) (*client.ServiceCollection, error) {
	if opts.Filters["name"] != serviceName {
		// TODO: Better error description
		return nil, errors.New(fmt.Sprintf("service should have received name: %s", serviceName))
	}
	// TODO: See if using real secondary launch config works
	slc := make(map[string]interface{})
	slc["ImageUuid"] = ""
	return &client.ServiceCollection{
		Data: []client.Service{
			client.Service{
				LaunchConfig: &client.LaunchConfig{
					ImageUuid: "",
				},
				SecondaryLaunchConfigs: []interface{}{
					slc,
					// client.SecondaryLaunchConfig{
					// 	ImageUuid: "",
					// },
				},
			},
		},
	}, nil
}

func (srv *SuccessfulService) ActionUpgrade(service *client.Service, upgrade *client.ServiceUpgrade) (*client.Service, error) {
	if upgrade.InServiceStrategy.SecondaryLaunchConfigs[0].(map[string]interface{})["imageUuid"] != fmt.Sprintf("docker:%s", codeTag) {
		return nil, errors.New("failed to upgrade service's code image")
	}
	return nil, nil
}

type UpgradeServiceService struct {
	SuccessfulService
}

func (srv *UpgradeServiceService) ActionUpgrade(service *client.Service, upgrade *client.ServiceUpgrade) (*client.Service, error) {
	if !validServiceVersionUpgrade(upgrade) {
		return nil, errors.New("upgrading service version failed")
	}
	return &client.Service{}, nil
}

type ServiceLikeName struct {
	SuccessfulService
}

func (srv *ServiceLikeName) List(opts *client.ListOpts) (*client.ServiceCollection, error) {
	if !validServiceLikeFilters(opts) {
		return nil, errors.New("invalid filters for service like query")
	}
	return &client.ServiceCollection{
		Data: []client.Service{
			client.Service{},
		},
	}, nil
}

func TestServiceByName(t *testing.T) {
	cli := Client{
		RancherClient: &client.RancherClient{
			Service: &FailedService{},
		},
	}

	_, err := cli.ServiceByName(serviceName)

	if err == nil {
		t.Errorf("retrieving service by name fails if rancher api call fails")
	}
}

func TestSuccessfulServiceByName(t *testing.T) {
	cli := Client{
		RancherClient: &client.RancherClient{
			Service: &SuccessfulService{},
		},
	}

	_, err := cli.ServiceByName(serviceName)

	if err != nil {
		t.Errorf("retrieving service by name fails if rancher api call fails")
	}
}

func TestUpgradeServiceCodeVersion(t *testing.T) {
	cli := Client{
		RancherClient: &client.RancherClient{
			Service: &SuccessfulService{},
		},
	}

	opts := UpgradeOpts{
		Service: serviceName,
		CodeTag: codeTag,
	}
	err := cli.UpgradeServiceCodeVersion(opts)

	if err != nil {
		t.Errorf("upgrading service with code version failed with %v", err)
	}
}

func TestFinishServiceUpgradeFailed(t *testing.T) {
	cli := Client{
		RancherClient: &client.RancherClient{
			Service: &FailedService{},
		},
	}

	_, err := cli.FinishServiceUpgrade(serviceName)

	if err == nil {
		t.Errorf("finishing service upgrade should have failed")
	}
}

func TestFinishServiceUpgrade(t *testing.T) {
	cli := Client{
		RancherClient: &client.RancherClient{
			Service: &SuccessfulService{},
		},
	}

	_, err := cli.FinishServiceUpgrade(serviceName)

	if err != nil {
		t.Errorf("finishing service upgrade failed with: %v", err)
	}
}

func TestServiceLikeName(t *testing.T) {
	cli := Client{
		RancherClient: &client.RancherClient{
			Service: &ServiceLikeName{},
		},
	}

	_, err := cli.ServiceLikeName(serviceName)

	if err != nil {
		t.Errorf("finishing service upgrade failed with: %v", err)
	}
}

func TestUpgradeServiceVersion(t *testing.T) {
	cli := Client{
		RancherClient: &client.RancherClient{
			Service: &UpgradeServiceService{},
		},
	}

	opts := UpgradeOpts{
		Service:    serviceName,
		RuntimeTag: codeTag,
	}
	err := cli.UpgradeServiceVersion(opts)

	if err != nil {
		t.Errorf("finishing service upgrade failed with: %v", err)
	}
}

func TestWaitTimesOutWhenUpgradeTakesTooLong(t *testing.T) {
	orig := upgradePollInterval
	upgradePollInterval = 10 * time.Second

	opts := UpgradeOpts{
		Interval: time.Millisecond,
	}
	cli := &Client{
		RancherClient: &client.RancherClient{
			Service: &UpgradeServiceService{},
		},
	}
	srv := &client.Service{
		Transitioning: "yes",
	}

	err := Wait(cli, srv, opts)

	if err == nil || err.Error() != "finishing upgrade timed out" {
		t.Errorf("service upgrade should have timed out")
	}

	upgradePollInterval = orig
}

func TestWaitReturnsNilWhenServiceIsNoLongerTransitioning(t *testing.T) {
	opts := UpgradeOpts{
		Interval: time.Millisecond,
	}
	cli := &Client{
		RancherClient: &client.RancherClient{
			Service: &UpgradeServiceService{},
		},
	}
	srv := &client.Service{
		Transitioning: "no",
	}

	err := Wait(cli, srv, opts)

	if err != nil {
		t.Errorf("wait should have exited cleanly")
	}
}

func TestUpgradeServiceWithName(t *testing.T) {

}

func validServiceLikeFilters(opts *client.ListOpts) bool {
	if opts.Filters["name_like"] != getServiceLikeQuery(serviceName) || opts.Filters["kind"] != SERVICE_TYPE_SERVICE {
		return false
	}
	return true
}

func validServiceVersionUpgrade(upgrade *client.ServiceUpgrade) bool {
	if upgrade.InServiceStrategy.LaunchConfig.ImageUuid != fmt.Sprintf("docker:%s", codeTag) {
		return false
	}
	return true
}
