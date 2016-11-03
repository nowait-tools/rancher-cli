package mocks

import (
	"errors"

	"github.com/rancher/go-rancher/client"
)

type FailedListEnvironmentOperations struct{}

func (env *FailedListEnvironmentOperations) List(opts *client.ListOpts) (*client.EnvironmentCollection, error) {
	return nil, errors.New("Failed to get environments for project")
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
	return nil, nil
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

func (env *FailedCreateEnvironmentOperations) List(opts *client.ListOpts) (*client.EnvironmentCollection, error) {
	return &client.EnvironmentCollection{
		Data: []client.Environment{
			client.Environment{},
		},
	}, nil
}

func (env *FailedCreateEnvironmentOperations) Create(opts *client.Environment) (*client.Environment, error) {
	return nil, errors.New("Failed to create environment")
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
			client.Environment{},
		},
	}, nil
}

func (env *SuccessfulEnvironmentOperations) Create(opts *client.Environment) (*client.Environment, error) {
	return opts, nil
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
