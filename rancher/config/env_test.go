package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rancher/go-rancher/client"
)

func TestGetEnvsReturnsAllEnvsInEnvFile(t *testing.T) {
	validator := EnvironmentValidator{
		EnvFilePath: "../../fixtures/.env",
	}
	envs, err := validator.getEnvs()

	if err != nil || len(envs) != 3 {
		t.Errorf("variable envs should contain 5 environment variables")
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		LaunchConfig         *client.LaunchConfig
		EnvironmentValidator EnvironmentValidator
		Error                error
	}{
		{
			LaunchConfig: getLaunchConfigWithEnvs("ENV_1", "ENV_2", "ENV_3"),
			EnvironmentValidator: EnvironmentValidator{
				EnvFilePath: "../../fixtures/.env",
			},
			Error: nil,
		},
		{

			LaunchConfig: getLaunchConfigWithEnvs("ENV_1"),
			EnvironmentValidator: EnvironmentValidator{
				EnvFilePath: "../../fixtures/.env",
			},
			Error: errors.New("env: missing ENV_2,ENV_3"),
		},
	}

	for _, test := range cases {

		err := test.EnvironmentValidator.Validate(test.LaunchConfig)

		fmt.Printf("error is %v\n", err)
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
