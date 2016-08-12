package cli

import "testing"

func TestExecuteOnFailedValidation(t *testing.T) {
	failedValidator := func(fileName string) bool {
		return false
	}
	cmd := &UpgradeServiceCommand{
		Validate: failedValidator,
	}

	err := cmd.Execute()

	if err == nil {
		t.Errorf("upgrade command should have returned validation error")
	}
}
