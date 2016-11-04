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

type FailedListEnvironmentOperations struct{}

func (env *FailedListEnvironmentOperations) List(opts *client.ListOpts) (*client.EnvironmentCollection, error) {
	return nil, ListEnvironmentsError
}
func (env *FailedListEnvironmentOperations) Create(opts *client.Environment) (*client.Environment, error) {
	return nil, nil
}
func (env *FailedListEnvironmentOperations) Update(existing *client.Environment, updates interface{}) (*client.Environment, error) {
	return nil, nil
}
func (env *FailedListEnvironmentOperations) ById(id string) (*client.Environment, error) {
	return nil, nil
}
func (env *FailedListEnvironmentOperations) Delete(container *client.Environment) error {
	return nil
}

func (env *FailedListEnvironmentOperations) ActionActivateservices(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionAddoutputs(*client.Environment, *client.AddOutputsInput) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionCancelrollback(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionCancelupgrade(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionCreate(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionDeactivateservices(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionError(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionExportconfig(*client.Environment, *client.ComposeConfigInput) (*client.ComposeConfig, error) {
	return nil, ActionExportconfigError
}

func (env *FailedListEnvironmentOperations) ActionFinishupgrade(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionRemove(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionRollback(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionUpdate(*client.Environment) (*client.Environment, error) {
	return nil, nil
}

func (env *FailedListEnvironmentOperations) ActionUpgrade(*client.Environment, *client.EnvironmentUpgrade) (*client.Environment, error) {
	return nil, nil
}

type SuccessfulEnvironmentOperations struct {
	FailedListEnvironmentOperations
}

type FailedCreateEnvironmentOperations struct {
	FailedListEnvironmentOperations
}

type FailedActionExportconfigEnvironmentOperations struct {
	FailedListEnvironmentOperations
}

func (env *FailedActionExportconfigEnvironmentOperations) ActionExportconfig(*client.Environment, *client.ComposeConfigInput) (*client.ComposeConfig, error) {
	return nil, ActionExportconfigError
}

func (env *FailedCreateEnvironmentOperations) List(opts *client.ListOpts) (*client.EnvironmentCollection, error) {
	return &client.EnvironmentCollection{
		Data: []client.Environment{
			client.Environment{},
		},
	}, nil
}

func (env *FailedActionExportconfigEnvironmentOperations) List(opts *client.ListOpts) (*client.EnvironmentCollection, error) {
	// Enforce that accountId_eq must be set to valid account string
	if validAccountid(opts.Filters) {
		// Should return EnvironmentCollection that has no results
		return &client.EnvironmentCollection{
			Data: []client.Environment{},
		}, nil
	}
	return &client.EnvironmentCollection{
		Data: []client.Environment{
			client.Environment{},
		},
	}, nil
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

func (env *SuccessfulEnvironmentOperations) List(opts *client.ListOpts) (*client.EnvironmentCollection, error) {
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
