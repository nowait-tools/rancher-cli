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
		Service              *client.Service
		EnvironmentValidator EnvironmentValidator
		Error                error
	}{
		{
			Service: &client.Service{
				LaunchConfig: getLaunchConfigWithEnvs("ENV_1", "ENV_2", "ENV_3"),
			},
			EnvironmentValidator: EnvironmentValidator{
				EnvFilePath: "../../fixtures/.env",
			},
			Error: nil,
		},
		{

			Service: &client.Service{
				LaunchConfig: getLaunchConfigWithEnvs("ENV_1"),
			},
			EnvironmentValidator: EnvironmentValidator{
				EnvFilePath: "../../fixtures/.env",
			},
			Error: errors.New("env: missing ENV_2,ENV_3"),
		},
	}

	for _, test := range cases {

		err := test.EnvironmentValidator.Validate(test.Service, UpgradeOpts{})

		fmt.Printf("error is %v\n", err)
		if errorsNotEqual(err, test.Error) {
			t.Errorf("validation error `%v` should match error expectation `%v`", err, test.Error)
		}
	}
}

func TestValidateEnvFlag(t *testing.T) {
	tests := []struct {
		Envs  []string
		Error bool
	}{
		{
			Envs: []string{
				"ENVIRONMENT=prod",
				"URL=https://test.salesforce.com/servlet/servlet.WebToLead?encoding=UTF-8",
			},
			Error: false,
		},
		{
			Envs: []string{
				"ENVIRONMENTprod",
			},
			Error: true,
		},
	}

	for _, test := range tests {
		result := ValidateEnvFlag(test.Envs)
		if test.Error {
			if result == nil {
				t.Errorf("validation for envs should have failed")
			}
		} else {
			if result != nil {
				t.Errorf("validation for envs should passed")
			}
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
