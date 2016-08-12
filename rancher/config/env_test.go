package config

import (
	"errors"
	"testing"

	"github.com/rancher/go-rancher/client"
)

func TestGetEnvsReturnsAllEnvsInEnvFile(t *testing.T) {
	envs, err := GetEnvs("../../fixtures/.env")

	if err != nil || len(envs) != 5 {
		t.Errorf("variable envs should contain 5 environment variables")
	}
}

func TestValidateEnvs(t *testing.T) {
	cases := []struct {
		LaunchConfig *client.LaunchConfig
		envs         []string
		Error        error
	}{
		{
			LaunchConfig: getLaunchConfigWithEnvs("ENV_1", "ENV_2", "ENV_3"),
			envs: []string{
				"ENV_1",
				"ENV_2",
				"ENV_3",
			},
			Error: nil,
		},
		{

			LaunchConfig: getLaunchConfigWithEnvs("ENV_1"),
			envs: []string{
				"ENV_1",
				"ENV_2",
				"ENV_3",
			},
			Error: errors.New("env: missing ENV_2,ENV_3"),
		},
	}

	for _, test := range cases {

		err := ValidateEnvs(test.envs, test.LaunchConfig)

		if errorsNotEqual(err, test.Error) {
			t.Errorf("validation error `%v` should match error expectation `%v`", err, test.Error)
		}
	}
}

func getLaunchConfigWithEnvs(envs ...string) *client.LaunchConfig {
	environment := make(map[string]interface{})

	for _, env := range envs {
		// Set each env as a key, value does not matter
		// since validation is based on key being present
		environment[env] = "Value"
	}
	return &client.LaunchConfig{
		Environment: environment,
	}
}

// Validates that errors check.  If errors are nil the comparing the pointers is enough.
// If the pointers are not nil, the value of the strings must be compared
func errorsNotEqual(actual, expected error) bool {
	return actual != expected && !(actual == nil || actual.Error() == expected.Error())
}
