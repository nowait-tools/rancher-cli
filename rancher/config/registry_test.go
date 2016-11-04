package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rancher/go-rancher/client"
)

var (
	sampleTags = []string{"1.0", "2.0"}

	failedToRetrieveTags = errors.New("failed to retrieve tags")
)

type NoopRegistryClient struct{}

func (client *NoopRegistryClient) Tags(repository string) (tags []string, err error) {
	return sampleTags, nil
}

type FailedRegistryClient struct{}

func (client *FailedRegistryClient) Tags(repository string) (tags []string, err error) {
	return nil, failedToRetrieveTags
}

type UnreliableRegsitryClient struct {
	count int
}

func (client *UnreliableRegsitryClient) Tags(repository string) (tags []string, err error) {
	client.count++
	if client.count == 2 {
		return nil, failedToRetrieveTags
	}

	return sampleTags, nil
}

func TestRegistryValidatorValidate(t *testing.T) {
	lc := make(map[string]interface{})
	lc["imageUuid"] = "docker:image/name:1.0"
	slcs := []interface{}{
		lc,
	}
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
				SecondaryLaunchConfigs: slcs,
			},
			Opts: UpgradeOpts{
				RuntimeTag: "image/name:2.0",
				CodeTag:    "image/name:2.0",
			},
			Validator: &RegistryValidator{
				RegistryClient: &NoopRegistryClient{},
			},
			Error:      nil,
			FailureMsg: "upgrade should be able to specify full name for code tag",
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
				SecondaryLaunchConfigs: slcs,
			},
			Opts: UpgradeOpts{
				CodeTag: "2.0",
			},
			Validator: &RegistryValidator{
				RegistryClient: &NoopRegistryClient{},
			},
			Error:      nil,
			FailureMsg: "upgrade should be able to specify the tag only for code tag",
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
				SecondaryLaunchConfigs: slcs,
			},
			Opts: UpgradeOpts{
				CodeTag: "3.0",
			},
			Validator: &RegistryValidator{
				RegistryClient: &NoopRegistryClient{},
			},
			Error:      ImageNotFound,
			FailureMsg: "should have received image not found for invalid tag for code tag",
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
				&FailedRegistryClient{},
			},
			Error:      failedToRetrieveTags,
			FailureMsg: "should have received failure on retrieving tags",
		},
		{
			Service: &client.Service{
				LaunchConfig: &client.LaunchConfig{
					ImageUuid: "docker:image/name:1.0",
				},
				SecondaryLaunchConfigs: slcs,
			},
			Opts: UpgradeOpts{
				CodeTag: "2.0",
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
			fmt.Printf("Test failed with error: %v", err)
		}
	}
}

func TestCachedRegsitryClient(t *testing.T) {
	repo := "repo"
	cache := make(map[string][]string)
	client := &CachedRegistryClient{
		Cache: cache,
		RegistryClient: &UnreliableRegsitryClient{
			count: 0,
		},
	}

	_, err := client.Tags(repo)

	if err != nil {
		t.Errorf("client should be retrieving the information for the cache")
	}

	_, err = client.Tags(repo)

	if err != nil {
		t.Errorf("client should be retrieving the information for the cache")
	}
}
