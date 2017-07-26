package rancher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nowait-tools/rancher-cli/rancher/config"
	"github.com/pkg/errors"
	"github.com/rancher/go-rancher/client"
)

const (
	SERVICE_TYPE_SERVICE = "service"
)

var (
	upgradePollInterval               = 150 * time.Millisecond
	environmentCloneSourceTargetError = errors.New("Could not find both source and target environments")
)

type Client struct {
	RancherClient *client.RancherClient
	Validators    []config.Validator
}

type UpgradeResult struct {
	Service *client.Service
	Error   error
}

// NewClient grabs config necessary and sets an inited client or returns an error
func NewClient(cattleURL string, cattleAccessKey string, cattleSecretKey string, envFile string) (*Client, error) {
	apiClient, err := client.NewRancherClient(&client.ClientOpts{
		Url:       cattleURL,
		AccessKey: cattleAccessKey,
		SecretKey: cattleSecretKey,
	})

	if err != nil {
		return nil, err
	}

	apiClient.Environment = &EnvironmentClient{
		apiClient.Environment,
		cattleAccessKey,
		cattleSecretKey,
		cattleURL,
	}

	registryValidator, err := config.NewRegistryValidator()

	if err != nil {
		return nil, err
	}

	if envFile != "" {
		return &Client{
			RancherClient: apiClient,
			Validators: []config.Validator{
				registryValidator,
				&config.EnvironmentValidator{
					EnvFilePath: envFile,
				},
			},
		}, nil

	} else {
		return &Client{
			RancherClient: apiClient,
			Validators: []config.Validator{
				registryValidator,
			},
		}, nil
	}
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

	if len(services.Data) != 1 {
		return nil, errors.New(fmt.Sprintf("failed to find service with name %s", name))
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

func (cli *Client) UpgradeService(opts config.UpgradeOpts) (*client.Service, error) {
	service, err := cli.ServiceByName(opts.Service)

	if err != nil {
		return service, err
	}

	if err = cli.ValidateService(service, opts); err != nil {
		return service, err
	}

	serviceUpgrade := UpdateLaunchConfig(service, opts)
	service, err = cli.RancherClient.Service.ActionUpgrade(service, serviceUpgrade)

	return service, err
}

// Clone the Rancher Project.  A project in rancher's api terms is equivalent to an environment.  And an environment is
// equivalent to a stack.
func (cli *Client) CloneProject(opts config.EnvUpgradeOpts) error {
	projects, err := cli.RancherClient.Project.List(&client.ListOpts{})

	if err != nil {
		return err
	}

	log.Debugf("Found %d projects", len(projects.Data))

	projMapping := struct {
		SourceProjectId string
		TargetProjectId string
	}{
		SourceProjectId: "",
		TargetProjectId: "",
	}
	// Use continue after finding matching project so that the opts.SourceProjectId != opts.TargetProjectId
	for _, project := range projects.Data {
		if project.Name == opts.SourceEnv {
			log.Debugf("Matched project %s with Id: %s", project.Name, project.Id)
			projMapping.SourceProjectId = project.Id
			continue
		}

		if project.Name == opts.TargetEnv {
			log.Debugf("Matched project %s with Id %s", project.Name, project.Id)
			projMapping.TargetProjectId = project.Id
			continue
		}
	}

	if projMapping.SourceProjectId == "" || projMapping.TargetProjectId == "" {
		return environmentCloneSourceTargetError
	}

	log.Debugf("Source project id %s target id %s", projMapping.SourceProjectId, projMapping.TargetProjectId)

	// Filter environments by correct project id and ensure they are active
	filters := make(map[string]interface{})
	filters["accountId_eq"] = projMapping.SourceProjectId
	filters["state"] = "active"
	envs, err := cli.RancherClient.Environment.List(&client.ListOpts{
		Filters: filters,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to find stacks for project")
	}

	log.Debugf("Found %d stacks, cloning into environment %s", len(envs.Data), opts.TargetEnv)

	for _, env := range envs.Data {
		composeConfig, err := cli.RancherClient.Environment.ActionExportconfig(&env, &client.ComposeConfigInput{})

		if err != nil {
			return err
		}

		_, err = cli.RancherClient.Environment.Create(&client.Environment{
			AccountId:      projMapping.TargetProjectId,
			DockerCompose:  composeConfig.DockerComposeConfig,
			RancherCompose: composeConfig.RancherComposeConfig,
			Name:           env.Name,
		})

		if err != nil {
			return err
		}
	}

	return err
}

// TODO: Simplify this method and test it
func (cli *Client) UpgradeServiceWithNameLike(opts config.UpgradeOpts) error {
	failed := false
	services, err := cli.ServiceLikeName(opts.ServiceLike)

	if err != nil {
		return err
	}

	serviceCount := len(services.Data)
	upgradeErrs := make(chan UpgradeResult, serviceCount)

	for _, service := range services.Data {
		go func(srv client.Service, opts config.UpgradeOpts) {
			opts.Service = srv.Name
			service, err := cli.UpgradeService(opts)

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
		}(service, opts)
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

func UpdateLaunchConfig(service *client.Service, opts config.UpgradeOpts) *client.ServiceUpgrade {
	inSrvStrat := &client.InServiceUpgradeStrategy{
		BatchSize:              1,
		IntervalMillis:         int64(opts.Interval) / (int64(math.Pow10(6))),
		StartFirst:             true,
		LaunchConfig:           service.LaunchConfig,
		SecondaryLaunchConfigs: service.SecondaryLaunchConfigs,
	}

	if opts.CodeTag != "" {

		lcImage := service.SecondaryLaunchConfigs[0].(map[string]interface{})["imageUuid"].(string)
		refs := strings.Split(opts.CodeTag, ":")
		image := ""
		switch len(refs) {
		case 1:
			first := strings.Index(lcImage, ":") + 1
			pos := strings.LastIndex(lcImage, ":") + 1
			image = lcImage[first:pos] + refs[0]
		case 2:
			image = opts.CodeTag
		}
		service.SecondaryLaunchConfigs[0].(map[string]interface{})["imageUuid"] = fmt.Sprintf("docker:%s", image)
		inSrvStrat.SecondaryLaunchConfigs = service.SecondaryLaunchConfigs
	}
	if opts.RuntimeTag != "" {

		refs := strings.Split(opts.RuntimeTag, ":")
		image := ""
		switch len(refs) {
		case 1:
			first := strings.Index(service.LaunchConfig.ImageUuid, ":") + 1
			pos := strings.LastIndex(service.LaunchConfig.ImageUuid, ":") + 1
			image = service.LaunchConfig.ImageUuid[first:pos] + refs[0]
		case 2:
			image = opts.RuntimeTag
			if opts.Repository != "" {
				image = fmt.Sprintf("%s/%s", opts.Repository, opts.RuntimeTag)
			}
		}
		service.LaunchConfig.ImageUuid = fmt.Sprintf("docker:%s", image)
		inSrvStrat.LaunchConfig = service.LaunchConfig
	}

	if len(opts.Envs) > 0 {
		for _, val := range opts.Envs {
			key, value := config.GetEnvValue(val)
			service.LaunchConfig.Environment[key] = value
			inSrvStrat.LaunchConfig = service.LaunchConfig
		}
	}

	return &client.ServiceUpgrade{
		Resource:          client.Resource{},
		InServiceStrategy: inSrvStrat,
	}
}

func Wait(cli *Client, srv *client.Service, opts config.UpgradeOpts) error {
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

func (cli *Client) ValidateService(service *client.Service, opts config.UpgradeOpts) error {
	for _, val := range cli.Validators {
		if err := val.Validate(service, opts); err != nil {
			return err
		}
	}

	return nil
}

func getServiceLikeQuery(serviceName string) string {
	return serviceName + "%"
}

type EnvironmentClient struct {
	client.EnvironmentOperations
	accessKey  string
	secretKey  string
	rancherUrl string
}

func (env *EnvironmentClient) Create(opts *client.Environment) (*client.Environment, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(*opts)
	url, err := url.Parse(env.rancherUrl)

	if err != nil {
		return nil, err
	}

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	createUrl := url.String() + "/projects/" + opts.AccountId + "/environments"

	req, err := http.NewRequest("POST", createUrl, &buf)
	req.SetBasicAuth(env.accessKey, env.secretKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 201 {
		return nil, errors.New(fmt.Sprintf("Creating environment returned %d response code", res.StatusCode))
	}

	_, err = ioutil.ReadAll(res.Body)

	return nil, err
}
