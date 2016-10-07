package rancher

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/nowait/rancher-cli/rancher/config"
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
			client.Service{
				LaunchConfig: &client.LaunchConfig{
					ImageUuid: "",
				},
			},
		},
	}, nil
}

type FailedValidator struct{}

func (val *FailedValidator) Validate(lc *client.LaunchConfig) error {
	return errors.New("validation has failed")
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

func TestUpgradeService(t *testing.T) {
	cli := Client{
		RancherClient: &client.RancherClient{
			Service: &UpgradeServiceService{},
		},
		Validator: &config.NoopValidator{},
	}

	opts := UpgradeOpts{
		Service:    serviceName,
		RuntimeTag: codeTag,
	}
	_, err := cli.UpgradeService(opts)

	if err != nil {
		t.Errorf("service upgrade failed with: %v", err)
	}
}

func TestUpgradeServiceFailsWhenValidationFails(t *testing.T) {
	cli := Client{
		RancherClient: &client.RancherClient{
			Service: &UpgradeServiceService{},
		},
		Validator: &FailedValidator{},
	}

	opts := UpgradeOpts{
		Service:    serviceName,
		RuntimeTag: codeTag,
	}
	_, err := cli.UpgradeService(opts)

	if err == nil {
		t.Errorf("service upgrade should have failed")
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

func TestUpdateLaunchConfig(t *testing.T) {
	slc := make(map[string]interface{})
	slc["imageUuid"] = ""
	srv := &client.Service{
		LaunchConfig: &client.LaunchConfig{
			ImageUuid: "",
		},
		SecondaryLaunchConfigs: []interface{}{
			slc,
		},
	}

	expectedSlc := make(map[string]interface{})
	expectedSlc["imageUuid"] = "docker:sample"

	tests := []struct {
		ExpectedServiceUpgrade *client.ServiceUpgrade
		Opts                   UpgradeOpts
	}{
		{
			ExpectedServiceUpgrade: &client.ServiceUpgrade{
				Resource: client.Resource{},
				InServiceStrategy: &client.InServiceUpgradeStrategy{
					BatchSize:      1,
					IntervalMillis: 10000,
					StartFirst:     true,
					LaunchConfig: &client.LaunchConfig{
						ImageUuid: "docker:sample",
					},
				},
			},
			Opts: UpgradeOpts{
				RuntimeTag: "sample",
			},
		},
		{
			ExpectedServiceUpgrade: &client.ServiceUpgrade{
				Resource: client.Resource{},
				InServiceStrategy: &client.InServiceUpgradeStrategy{
					BatchSize:      1,
					IntervalMillis: 10000,
					StartFirst:     true,
					SecondaryLaunchConfigs: []interface{}{
						expectedSlc,
					},
				},
			},
			Opts: UpgradeOpts{
				CodeTag: "sample",
			},
		},
		{
			ExpectedServiceUpgrade: &client.ServiceUpgrade{
				Resource: client.Resource{},
				InServiceStrategy: &client.InServiceUpgradeStrategy{
					BatchSize:      1,
					IntervalMillis: 10000,
					StartFirst:     true,
					LaunchConfig: &client.LaunchConfig{
						ImageUuid: "docker:sample",
					},
					SecondaryLaunchConfigs: []interface{}{
						expectedSlc,
					},
				},
			},
			Opts: UpgradeOpts{
				RuntimeTag: "sample",
				CodeTag:    "sample",
			},
		},
	}

	for _, test := range tests {
		actual := UpdateLaunchConfig(srv, test.Opts)

		if !reflect.DeepEqual(actual, test.ExpectedServiceUpgrade) {
			t.Errorf("failure")
			fmt.Printf("%v", pretty.Diff(actual, test.ExpectedServiceUpgrade))
		}
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
