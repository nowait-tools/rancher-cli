package rancher

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/imdario/mergo"
	"github.com/kr/pretty"
	"github.com/nowait/rancher-cli/rancher/config"
	"github.com/nowait/rancher-cli/rancher/mocks"
	"github.com/rancher/go-rancher/client"
)

var serviceName = "name"
var codeTag = "image-name:1.0"
var defaultImageUuid = "docker:runtime/image:1.0"
var defaultSlcImageUuid = "docker:code/image:1.0"

var upgradedUpgradeOptsImage = "runtime/image:2.0"
var upgradedCodeOpts = "code/image:2.0"

var upgradedImageUuid = "docker:runtime/image:2.0"
var upgradedSlcImageUuid = "docker:code/image:2.0"

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
					ImageUuid: defaultImageUuid,
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
	return dummyService(), nil
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

func (val *FailedValidator) Validate(service *client.Service, opts config.UpgradeOpts) error {
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
		Validators: []config.Validator{
			&config.NoopValidator{},
		},
	}

	opts := config.UpgradeOpts{
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
		Validators: []config.Validator{
			&FailedValidator{},
		},
	}

	opts := config.UpgradeOpts{
		Service:    serviceName,
		RuntimeTag: codeTag,
	}
	_, err := cli.UpgradeService(opts)

	if err == nil {
		t.Errorf("service upgrade should have failed")
	}
}

func TestUpgradeServiceFailsWhenSingleValidatorFails(t *testing.T) {
	cli := Client{
		RancherClient: &client.RancherClient{
			Service: &UpgradeServiceService{},
		},
		Validators: []config.Validator{
			&config.NoopValidator{},
			&FailedValidator{},
		},
	}

	opts := config.UpgradeOpts{
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

	opts := config.UpgradeOpts{
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
	opts := config.UpgradeOpts{
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
	expectedSlc := make(map[string]interface{})
	expectedSlc["imageUuid"] = "docker:sample"
	var expectedInterval int64 = 5000
	interval := 5 * time.Second

	expectedEnvs := make(map[string]interface{})
	expectedEnvs["ENVIRONMENT"] = "prod"

	tests := []struct {
		ExpectedServiceUpgrade *client.ServiceUpgrade
		Opts                   config.UpgradeOpts
	}{
		{
			ExpectedServiceUpgrade: expectedServiceUpgrade(serviceUpgradeOverrides{
				ImageUuid: upgradedImageUuid,
				Interval:  expectedInterval,
			}),
			Opts: config.UpgradeOpts{
				RuntimeTag: upgradedUpgradeOptsImage,
				Interval:   interval,
			},
		},
		{
			ExpectedServiceUpgrade: expectedServiceUpgrade(serviceUpgradeOverrides{
				SlcImageUuid: upgradedSlcImageUuid,
				Interval:     expectedInterval,
			}),
			Opts: config.UpgradeOpts{
				CodeTag:  upgradedCodeOpts,
				Interval: interval,
			},
		},
		{
			ExpectedServiceUpgrade: expectedServiceUpgrade(serviceUpgradeOverrides{
				ImageUuid:    upgradedImageUuid,
				SlcImageUuid: upgradedSlcImageUuid,
				Interval:     expectedInterval,
			}),
			Opts: config.UpgradeOpts{
				RuntimeTag: "2.0",
				CodeTag:    "2.0",
				Interval:   interval,
			},
		},
		{
			ExpectedServiceUpgrade: expectedServiceUpgrade(serviceUpgradeOverrides{
				Environment: []string{
					"ENVIRONMENT=prod",
				},
				Interval: expectedInterval,
			}),
			Opts: config.UpgradeOpts{
				Envs: []string{
					"ENVIRONMENT=prod",
				},
				Interval: interval,
			},
		},
	}

	for index, test := range tests {
		actual := UpdateLaunchConfig(dummyService(), test.Opts)

		if !reflect.DeepEqual(actual, test.ExpectedServiceUpgrade) {
			t.Errorf("failure for test case %d", index)
			fmt.Printf("actual: %#v \n\n expected: %#v \n\n Diff: %v", actual, test.ExpectedServiceUpgrade, pretty.Diff(actual, test.ExpectedServiceUpgrade))
		}
	}
}

func TestCloneProject(t *testing.T) {
	// TODO: Update tests to verify that the correct error case is happening, could be somewhat dangeous to refactor the test cases without having some type of verifcation that the error cases are the correct one.

	tests := []struct {
		Description string
		Client      Client
		Opts        config.EnvUpgradeOpts
		ShouldFail  bool
	}{
		{
			Description: "When retrieving the projects from Rancher fails",
			Client: Client{
				RancherClient: &client.RancherClient{

					Environment: &mocks.SuccessfulEnvironmentOperations{},
					Project:     &mocks.FailedProjectOperations{},
				},
			},
			Opts: config.EnvUpgradeOpts{
				SourceEnv: mocks.ProjectOneName,
				TargetEnv: mocks.ProjectTwoName,
			},
			ShouldFail: true,
		},
		{
			Description: "When source environment not found in Rancher",
			Client: Client{
				RancherClient: &client.RancherClient{

					Environment: &mocks.SuccessfulEnvironmentOperations{},
					Project:     &mocks.PartialProjectOperations{},
				},
			},
			Opts: config.EnvUpgradeOpts{
				SourceEnv: "not found",
				TargetEnv: mocks.ProjectTwoName,
			},
			ShouldFail: true,
		},
		{
			Description: "When target environment not found in Rancher",
			Client: Client{
				RancherClient: &client.RancherClient{

					Environment: &mocks.SuccessfulEnvironmentOperations{},
					Project:     &mocks.PartialProjectOperations{},
				},
			},
			Opts: config.EnvUpgradeOpts{
				SourceEnv: mocks.ProjectOneName,
				TargetEnv: "not found",
			},
			ShouldFail: true,
		},
		{
			Description: "When using same source and target environments",
			Client: Client{
				RancherClient: &client.RancherClient{

					Environment: &mocks.SuccessfulEnvironmentOperations{},
					Project:     &mocks.PartialProjectOperations{},
				},
			},
			Opts: config.EnvUpgradeOpts{
				SourceEnv: mocks.ProjectOneName,
				TargetEnv: mocks.ProjectOneName,
			},
			ShouldFail: true,
		},
		{
			Description: "When listing the environments in Rancher fails",
			Client: Client{
				RancherClient: &client.RancherClient{

					Environment: &mocks.FailedListEnvironmentOperations{},
					Project:     &mocks.SuccessfulProjectOperations{},
				},
			},
			Opts: config.EnvUpgradeOpts{
				SourceEnv: mocks.ProjectOneName,
				TargetEnv: mocks.ProjectTwoName,
			},
			ShouldFail: true,
		},
		{
			Description: "When creating the Rancher environment fails",
			Client: Client{
				RancherClient: &client.RancherClient{

					Environment: &mocks.FailedCreateEnvironmentOperations{},
					Project:     &mocks.SuccessfulProjectOperations{},
				},
			},
			Opts: config.EnvUpgradeOpts{
				SourceEnv: mocks.ProjectOneName,
				TargetEnv: mocks.ProjectTwoName,
			},
			ShouldFail: true,
		},
		{
			Description: "Successful clone of environment",
			Client: Client{
				RancherClient: &client.RancherClient{

					Environment: &mocks.SuccessfulEnvironmentOperations{},
					Project:     &mocks.SuccessfulProjectOperations{},
				},
			},
			Opts: config.EnvUpgradeOpts{
				SourceEnv: mocks.ProjectOneName,
				TargetEnv: mocks.ProjectTwoName,
			},
			ShouldFail: false,
		},
	}

	for index, test := range tests {

		err := test.Client.CloneProject(test.Opts)

		if test.ShouldFail && err == nil {
			t.Errorf("Test case %d failed, expected error but received nil", index)
		}

		if !test.ShouldFail && err != nil {
			t.Errorf("Test case %d failed, received error %v", index, err)
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

func dummyService() *client.Service {
	slc := make(map[string]interface{})
	slc["imageUuid"] = defaultSlcImageUuid

	envs := make(map[string]interface{})
	return &client.Service{
		LaunchConfig: &client.LaunchConfig{
			ImageUuid:   defaultImageUuid,
			Environment: envs,
		},
		SecondaryLaunchConfigs: []interface{}{
			slc,
		},
	}
}

type serviceUpgradeOverrides struct {
	ImageUuid    string
	Environment  []string
	Interval     int64
	SlcImageUuid string
}

func (overrides serviceUpgradeOverrides) environmentVariables() map[string]interface{} {
	environment := make(map[string]interface{})
	for _, env := range overrides.Environment {
		splits := strings.Split(env, "=")
		environment[splits[0]] = splits[1]
	}
	return environment
}

func expectedServiceUpgrade(overrides serviceUpgradeOverrides) *client.ServiceUpgrade {
	defaults := serviceUpgradeOverrides{
		ImageUuid:    defaultImageUuid,
		Environment:  []string{},
		Interval:     10000,
		SlcImageUuid: defaultSlcImageUuid,
	}

	if err := mergo.Merge(&overrides, defaults); err != nil {
		fmt.Printf("failed to merge structs: %v", err)
	}

	expectedSlc := make(map[string]interface{})
	expectedSlc["imageUuid"] = overrides.SlcImageUuid
	return &client.ServiceUpgrade{
		Resource: client.Resource{},
		InServiceStrategy: &client.InServiceUpgradeStrategy{
			BatchSize:      1,
			IntervalMillis: overrides.Interval,
			StartFirst:     true,
			LaunchConfig: &client.LaunchConfig{
				ImageUuid:   overrides.ImageUuid,
				Environment: overrides.environmentVariables(),
			},
			SecondaryLaunchConfigs: []interface{}{
				expectedSlc,
			},
		},
	}
}
