package config

import (
	"errors"
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

func (val *EnvironmentValidator) Validate(service *client.Service, opts UpgradeOpts) error {
	missing := []string{}
	envs, err := val.getEnvs()
	lc := service.LaunchConfig

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

// Validate that the strings contained in a slice match the form
// key=value.
func ValidateEnvFlag(envs []string) error {
	for _, env := range envs {
		if pieces := strings.SplitN(env, "=", 2); len(pieces) != 2 {
			return errors.New(fmt.Sprintf("invalid env: %v\n expected key value pair in the form key=value", env))
		}
	}
	return nil
}

// Function to return the key and value of the form key=value.
// This method assumes it is already a valid environment variable
// string.
func GetEnvValue(env string) (key, value string) {
	pieces := strings.SplitN(env, "=", 2)
	return pieces[0], pieces[1]
}
