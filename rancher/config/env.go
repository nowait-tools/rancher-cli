package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rancher/go-rancher/client"
)

func GetEnvs(filenames ...string) ([]string, error) {
	envMap, err := godotenv.Read(filenames...)

	if err != nil {
		return nil, err
	}
	keys := []string{}
	for k := range envMap {
		keys = append(keys, k)
	}
	return keys, nil
}

func ValidateEnvs(envs []string, lc *client.LaunchConfig) error {
	missing := []string{}
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
