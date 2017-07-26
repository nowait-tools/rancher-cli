package config

import (
	"time"

	"github.com/rancher/go-rancher/client"
)

type Validator interface {
	Validate(service *client.Service, opts UpgradeOpts) error
}

type NoopValidator struct{}

func (val *NoopValidator) Validate(service *client.Service, opts UpgradeOpts) error {
	return nil
}

type UpgradeOpts struct {
	Envs        []string
	Wait        bool
	ServiceLike string
	Service     string
	CodeTag     string
	RuntimeTag  string
	Repository  string
	Interval    time.Duration
}

type EnvUpgradeOpts struct {
	SourceEnv string
	TargetEnv string
}
