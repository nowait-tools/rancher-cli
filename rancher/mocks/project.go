package mocks

import (
	"errors"

	"github.com/rancher/go-rancher/client"
)

var ProjectOneName = "one"
var ProjectTwoName = "two"
var ProjectOne = &client.Project{Name: ProjectOneName}
var ProjectTwo = &client.Project{Name: ProjectTwoName}

type FailedProjectOperations struct{}

func (proj *FailedProjectOperations) List(opts *client.ListOpts) (*client.ProjectCollection, error) {
	return nil, errors.New("Failed to list projects from Rancher api")
}

func (proj *FailedProjectOperations) Create(opts *client.Project) (*client.Project, error) {
	return nil, nil
}

func (proj *FailedProjectOperations) Update(existing *client.Project, updates interface{}) (*client.Project, error) {
	return nil, nil
}

func (proj *FailedProjectOperations) ById(id string) (*client.Project, error) {
	return nil, nil
}

func (proj *FailedProjectOperations) Delete(container *client.Project) error {
	return nil
}

func (proj *FailedProjectOperations) ActionActivate(*client.Project) (*client.Account, error) {
	return nil, nil
}

func (proj *FailedProjectOperations) ActionCreate(*client.Project) (*client.Account, error) {
	return nil, nil
}

func (proj *FailedProjectOperations) ActionDeactivate(*client.Project) (*client.Account, error) {
	return nil, nil
}

func (proj *FailedProjectOperations) ActionPurge(*client.Project) (*client.Account, error) {
	return nil, nil
}

func (proj *FailedProjectOperations) ActionRemove(*client.Project) (*client.Account, error) {
	return nil, nil
}

func (proj *FailedProjectOperations) ActionRestore(*client.Project) (*client.Account, error) {
	return nil, nil
}

func (proj *FailedProjectOperations) ActionSetmembers(*client.Project, *client.SetProjectMembersInput) (*client.SetProjectMembersInput, error) {
	return nil, nil
}

func (proj *FailedProjectOperations) ActionUpdate(*client.Project) (*client.Account, error) {
	return nil, nil
}

type PartialProjectOperations struct {
	FailedProjectOperations
}

func (proj *PartialProjectOperations) List(opts *client.ListOpts) (*client.ProjectCollection, error) {
	return &client.ProjectCollection{
		Data: []client.Project{
			*ProjectOne,
		},
	}, nil
}

type SuccessfulProjectOperations struct {
	FailedProjectOperations
}

func (proj *SuccessfulProjectOperations) List(opts *client.ListOpts) (*client.ProjectCollection, error) {
	return &client.ProjectCollection{
		Data: []client.Project{
			*ProjectOne,
			*ProjectTwo,
		},
	}, nil
}
