package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/rancher/go-rancher/client"
)

var (
	ImageNotFound = errors.New("image not found")
	username      = os.Getenv("DOCKER_REGISTRY_USERNAME")
	password      = os.Getenv("DOCKER_REGISTRY_PASSWORD")
	registryUrl   = os.Getenv("DOCKER_REGISTRY_URL")
)

type RegistryValidator struct {
	RegistryClient RegistryClient
}

type RegistryClient interface {
	Tags(repository string) (tags []string, err error)
}

func NewRegistryValidator() (*RegistryValidator, error) {
	client, err := registry.New(registryUrl, username, password)

	if err != nil {
		return nil, err
	}

	return &RegistryValidator{
		RegistryClient: client,
	}, nil
}

func (val *RegistryValidator) Validate(service *client.Service, opts UpgradeOpts) error {
	// Verify that the image name provided is valid and that it exists in the registry
	lc := service.LaunchConfig
	upgradeCode := opts.RuntimeTag
	// Check if image is valid in the context of docker's rules
	ref, err := validateImageName(upgradeCode)

	if err != nil {
		return err
	}

	repo := ""
	expectedTag := ""
	switch t := ref.(type) {
	case reference.NamedTagged:
		repo = t.Name()
		expectedTag = t.Tag()
	case reference.Reference:
		repo = imageUuidToRepository(lc.ImageUuid)
		expectedTag = t.String()
	default:
		panic(fmt.Sprintf("unsupported reference type %v", ref))
	}

	tags, err := val.RegistryClient.Tags(repo)

	if err != nil {
		return err
	}

	for _, tag := range tags {
		if tag == expectedTag {
			return nil
		}
	}

	return ImageNotFound
}

func validateImageName(image string) (reference.Reference, error) {
	ref, err := reference.Parse(image)

	if err != nil {
		return nil, err
	}

	return ref, nil
}

// Given a LaunchConfig.ImageUuid of the form docker:image/name:tag
// return just the image/name
func imageUuidToRepository(imageUuid string) string {
	return strings.Split(imageUuid, ":")[1]
}
