package config

import (
	"errors"
	"testing"

	"github.com/rancher/go-rancher/client"
)

type NoopRegistryClient struct{}

func (client *NoopRegistryClient) Tags(repository string) (tags []string, err error) {
	return []string{
		"1.0",
		"2.0",
	}, nil
}

type FailedRegistryClient struct{}

var failedToRetrieveTags = errors.New("failed to retrieve tags")

func (client *FailedRegistryClient) Tags(repository string) (tags []string, err error) {
	return nil, failedToRetrieveTags
}

func TestRegistryValidatorValidate(t *testing.T) {
	tests := []struct {
		Service    *client.Service
		Opts       UpgradeOpts
		Validator  *RegistryValidator
		Error      error
		FailureMsg string
	}{
		{
			Service: &client.Service{
				LaunchConfig: &client.LaunchConfig{
					ImageUuid: "docker:image/name:1.0",
				},
			},
			Opts: UpgradeOpts{
				RuntimeTag: "image/name:2.0",
			},
			Validator: &RegistryValidator{
				RegistryClient: &NoopRegistryClient{},
			},
			Error:      nil,
			FailureMsg: "upgrade should be able to specify full name for runtime tag",
		},
		{
			Service: &client.Service{
				LaunchConfig: &client.LaunchConfig{
					ImageUuid: "docker:image/name:1.0",
				},
			},
			Opts: UpgradeOpts{
				RuntimeTag: "2.0",
			},
			Validator: &RegistryValidator{
				&NoopRegistryClient{},
			},
			Error:      nil,
			FailureMsg: "upgrade should be able to specify the tag only for runtime tag",
		},
		{
			Service: &client.Service{
				LaunchConfig: &client.LaunchConfig{
					ImageUuid: "docker:image/name:1.0",
				},
			},
			Opts: UpgradeOpts{
				RuntimeTag: "3.0",
			},
			Validator: &RegistryValidator{
				&NoopRegistryClient{},
			},
			Error:      ImageNotFound,
			FailureMsg: "should have received image not found for invalid tag for runtime tag",
		},
		{
			Service: &client.Service{
				LaunchConfig: &client.LaunchConfig{
					ImageUuid: "docker:image/name:1.0",
				},
			},
			Opts: UpgradeOpts{
				RuntimeTag: "3.0",
			},
			Validator: &RegistryValidator{
				&FailedRegistryClient{},
			},
			Error:      failedToRetrieveTags,
			FailureMsg: "should have received failure on retrieving tags",
		},
	}

	for _, test := range tests {

		err := test.Validator.Validate(test.Service, test.Opts)

		if err != test.Error {
			t.Errorf(test.FailureMsg)
		}
	}
}

func TestValidateImageName(t *testing.T) {
	tests := []struct {
		Image      string
		ShouldFail bool
	}{
		{
			Image:      "valid/image:0.2",
			ShouldFail: false,
		},
		{
			// Tags cannot start with a period
			Image:      "invalid/image:.0.2",
			ShouldFail: true,
		},
		{
			// Tag names should also be accepted
			Image:      "0.2.0",
			ShouldFail: false,
		},
	}

	for _, test := range tests {

		_, err := validateImageName(test.Image)

		if (!test.ShouldFail && err != nil) || (test.ShouldFail && err == nil) {
			t.Errorf("test failed")
		}
	}
}
