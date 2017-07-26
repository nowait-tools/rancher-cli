package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/rancher/go-rancher/client"
	"log"
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

type cachedRegistryClient struct {
	cache          map[string][]string
	registryClient RegistryClient
}

func (cache *cachedRegistryClient) Tags(repository string) (tags []string, err error) {
	if _, ok := cache.cache[repository]; ok {
		return cache.cache[repository], nil
	}

	tags, err = cache.registryClient.Tags(repository)

	cache.cache[repository] = tags

	return
}

func newCachedRegistryClient(registryUrl, username, password string) (RegistryClient, error) {
	cache := make(map[string][]string)
	client, err := registry.New(registryUrl, username, password)
	return &cachedRegistryClient{
		registryClient: client,
		cache:          cache,
	}, err
}

type image struct {
	launchConfigImage string
	upgradeImage      string
}

func NewRegistryValidator() (*RegistryValidator, error) {
	log.Println("registryUrl:", registryUrl)
	log.Println("username:", username)
	client, err := newCachedRegistryClient(registryUrl, username, password)

	if err != nil {
		return nil, err
	}

	return &RegistryValidator{
		RegistryClient: client,
	}, nil
}

func (val *RegistryValidator) Validate(service *client.Service, opts UpgradeOpts) error {
	// Verify that the image name provided is valid and that it exists in the registry
	images := []image{}
	if opts.RuntimeTag != "" {

		images = append(images, image{
			upgradeImage:      opts.RuntimeTag,
			launchConfigImage: service.LaunchConfig.ImageUuid,
		})
	}
	if opts.CodeTag != "" {

		images = append(images, image{
			launchConfigImage: service.SecondaryLaunchConfigs[0].(map[string]interface{})["imageUuid"].(string),
			upgradeImage:      opts.CodeTag,
		})
	}

	return val.imageExistsInRegistry(images)
}

func (val *RegistryValidator) imageExistsInRegistry(images []image) error {
	for _, image := range images {
		ref, err := reference.Parse(image.upgradeImage)

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
			repo = imageUuidToRepository(image.launchConfigImage)
			expectedTag = t.String()
		default:
			panic(fmt.Sprintf("unsupported reference type %v", ref))
		}

		tags, err := val.RegistryClient.Tags(repo)

		if err != nil {
			return err
		}

		if containsTag(expectedTag, tags) {
			continue
		}

		return ImageNotFound

	}
	return nil
}

// Given a LaunchConfig.ImageUuid of the form docker:image/name:tag
// return just the image/name
func imageUuidToRepository(imageUuid string) string {
	return strings.Split(imageUuid, ":")[1]
}

func containsTag(expectedTag string, tags []string) bool {
	for _, tag := range tags {
		if tag == expectedTag {
			return true
		}
	}
	return false
}
