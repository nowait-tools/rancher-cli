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

type FailedService struct {
	ListOpts *client.ListOpts
}

func (srv *FailedService) List(opts *client.ListOpts) (*client.ServiceCollection, error) {
	return nil, errors.New("retrieving services failed with error")
}
func (srv *FailedService) Create(opts *client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) Update(existing *client.Service, updates interface{}) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ById(id string) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) Delete(container *client.Service) error {
	return nil
}
func (srv *FailedService) ActionActivate(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionAddservicelink(*client.Service, *client.AddRemoveServiceLinkInput) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionCancelrollback(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionCancelupgrade(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionCreate(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionDeactivate(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionFinishupgrade(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionRemove(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionRemoveservicelink(*client.Service, *client.AddRemoveServiceLinkInput) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionRestart(*client.Service, *client.ServiceRestart) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionRollback(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *FailedService) ActionSetservicelinks(*client.Service, *client.SetServiceLinksInput) (*client.Service, error) {
	return nil, nil
}

func (srv *FailedService) ActionUpdate(*client.Service) (*client.Service, error) {
	return nil, nil
}

func (srv *FailedService) ActionUpgrade(*client.Service, *client.ServiceUpgrade) (*client.Service, error) {
	return nil, nil
}

type SuccessfulService struct {
	ListOpts *client.ListOpts
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

func (srv *SuccessfulService) Create(opts *client.Service) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) Update(existing *client.Service, updates interface{}) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ById(id string) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) Delete(container *client.Service) error {
	return nil
}

func (srv *SuccessfulService) ActionActivate(*client.Service) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ActionAddservicelink(*client.Service, *client.AddRemoveServiceLinkInput) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ActionCancelrollback(*client.Service) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ActionCancelupgrade(*client.Service) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ActionCreate(*client.Service) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ActionDeactivate(*client.Service) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ActionFinishupgrade(*client.Service) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ActionRemove(*client.Service) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ActionRemoveservicelink(*client.Service, *client.AddRemoveServiceLinkInput) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ActionRestart(*client.Service, *client.ServiceRestart) (*client.Service, error) {
	return nil, nil
}
func (srv *SuccessfulService) ActionRollback(*client.Service) (*client.Service, error) {
	return nil, nil
}
func (srv *SuccessfulService) ActionSetservicelinks(*client.Service, *client.SetServiceLinksInput) (*client.Service, error) {
	return nil, nil
}

func (srv *SuccessfulService) ActionUpdate(*client.Service) (*client.Service, error) {
	return nil, nil
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
		&client.RancherClient{
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
		&client.RancherClient{
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
		&client.RancherClient{
			Service: &SuccessfulService{},
		},
	}

	err := cli.UpgradeServiceCodeVersion(serviceName, codeTag)

	if err != nil {
		t.Errorf("upgrading service with code version failed with %v", err)
	}
}

func TestFinishServiceUpgradeFailed(t *testing.T) {
	cli := Client{
		&client.RancherClient{
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
		&client.RancherClient{
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
		&client.RancherClient{
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
		&client.RancherClient{
			Service: &UpgradeServiceService{},
		},
	}

	err := cli.UpgradeServiceVersion(serviceName, codeTag)

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
		&client.RancherClient{
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
		&client.RancherClient{
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
