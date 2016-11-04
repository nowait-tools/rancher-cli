package mocks

import (
	"errors"

	"github.com/rancher/go-rancher/client"
)

var (
	ActionExportconfigError = errors.New("Error exporting compose config")
	ListEnvironmentsError   = errors.New("Failed to find environments for prjoect")
	CreateEnvironmentError  = errors.New("Failed to create environments for prjoect")
)

// TODO: Need to clean up this mess
type NoopEnvironmentOperations struct{}

func (env *NoopEnvironmentOperations) Create(opts *client.Environment) (*client.Environment, error) {
	return nil, nil
}
func (env *NoopEnvironmentOperations) Update(existing *client.Environment, updates interface{}) (*client.Environment, error) {
	return nil, nil
}
func (env *NoopEnvironmentOperations) ById(id string) (*client.Environment, error) {
	return nil, nil
}
func (env *NoopEnvironmentOperations) Delete(container *client.Environment) error {
	return nil
}

func (env *NoopEnvironmentOperations) ActionActivateservices(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionAddoutputs(*client.Environment, *client.AddOutputsInput) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionCancelrollback(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionCancelupgrade(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionCreate(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionDeactivateservices(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionError(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionExportconfig(*client.Environment, *client.ComposeConfigInput) (*client.ComposeConfig, error) {
	return nil, ActionExportconfigError
}

func (env *NoopEnvironmentOperations) ActionFinishupgrade(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionRemove(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionRollback(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionUpdate(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *NoopEnvironmentOperations) ActionUpgrade(*client.Environment, *client.EnvironmentUpgrade) (*client.Environment, error) {
	return nil, nil
}

type FailedListEnvironmentOperations struct {
	NoopEnvironmentOperations
}

func (env *FailedListEnvironmentOperations) List(opts *client.ListOpts) (*client.EnvironmentCollection, error) {
	return nil, ListEnvironmentsError
}

type SuccessfulEnvironmentOperations struct {
	NoopEnvironmentOperations
	SuccessfulList
}

type FailedCreateEnvironmentOperations struct {
	NoopEnvironmentOperations
	SuccessfulList
}

type FailedActionExportconfigEnvironmentOperations struct {
	SuccessfulList
	NoopEnvironmentOperations
}

func (env *FailedActionExportconfigEnvironmentOperations) ActionExportconfig(*client.Environment, *client.ComposeConfigInput) (*client.ComposeConfig, error) {
	return nil, ActionExportconfigError
}

func (env *FailedCreateEnvironmentOperations) ActionExportconfig(*client.Environment, *client.ComposeConfigInput) (*client.ComposeConfig, error) {
	return &client.ComposeConfig{
		DockerComposeConfig:  "",
		RancherComposeConfig: "",
	}, nil
}

func (env *FailedCreateEnvironmentOperations) Create(opts *client.Environment) (*client.Environment, error) {
	return nil, CreateEnvironmentError
}

type SuccessfulList struct{}

func (env *SuccessfulList) List(opts *client.ListOpts) (*client.EnvironmentCollection, error) {
	// Enforce that accountId_eq must be set to valid account string
	if validAccountid(opts.Filters) {
		// Should return EnvironmentCollection that has no results
		return &client.EnvironmentCollection{
			Data: []client.Environment{},
		}, nil
	}
	return &client.EnvironmentCollection{
		Data: []client.Environment{
			client.Environment{
				Name: "name",
			},
		},
	}, nil
}

func (env *SuccessfulEnvironmentOperations) Create(opts *client.Environment) (*client.Environment, error) {
	if opts.AccountId == "" || opts.DockerCompose == "" || opts.RancherCompose == "" || opts.Name == "" {
		return nil, errors.New("Failed to create environment: accountId, docker compose, rancher compose or name empty string")
	}
	return opts, nil
}

func (env *SuccessfulEnvironmentOperations) ActionExportconfig(*client.Environment, *client.ComposeConfigInput) (*client.ComposeConfig, error) {
	return &client.ComposeConfig{
		DockerComposeConfig:  "docker compose",
		RancherComposeConfig: "rancher compose",
	}, nil
}

func validAccountid(filters map[string]interface{}) bool {
	acctId, ok := filters["accountId_eq"].(string)

	if !ok {
		acctId = filters["accountId"].(string)
	}

	if len(acctId) < 2 {
		return false
	}

	// Rancher always prefixes account ids with 1a
	data := []byte(acctId)
	if data[0] != '1' || data[1] != 'a' {
		return false
	}

	return true
}
