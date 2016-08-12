package cli

type UpgradeServiceCommand struct {
	Validator UpgradeValidator
}

type UpgradeValidator struct {
}

func (cmd *UpgradeServiceCommand) Execute() error {

	return nil
}
