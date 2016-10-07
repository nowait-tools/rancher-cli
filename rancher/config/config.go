package config

import "github.com/rancher/go-rancher/client"

type Validator interface {
	Validate(lc *client.LaunchConfig) error
}

type NoopValidator struct{}

func (val *NoopValidator) Validate(lc *client.LaunchConfig) error {
	return nil
}
