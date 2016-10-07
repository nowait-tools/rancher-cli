package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rancher/go-rancher/client"
)

type EnvironmentValidator struct {
	EnvFilePath string
}

func (val *EnvironmentValidator) getEnvs() ([]string, error) {
	envMap, err := godotenv.Read(val.EnvFilePath)

	if err != nil {
		return nil, err
	}
	keys := []string{}
	for k := range envMap {
		keys = append(keys, k)
	}
	return keys, nil
}

func (val *EnvironmentValidator) Validate(lc *client.LaunchConfig) error {
	missing := []string{}
	envs, err := val.getEnvs()

	if err != nil {
		return err
	}
	for _, env := range envs {
		_, ok := lc.Environment[env]

		if !ok {
			missing = append(missing, env)
		}
	}

	if len(missing) != 0 {
		return fmt.Errorf("env: missing %s", strings.Join(missing, ","))
	}
	return nil
}
