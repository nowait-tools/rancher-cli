package config

import "github.com/rancher/go-rancher/client"

type Validator interface {
	Validate(lc *client.LaunchConfig) error
}
