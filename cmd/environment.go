package cmd

import (
	"github.com/nowait-tools/rancher-cli/rancher"
	"github.com/nowait-tools/rancher-cli/rancher/config"
	"github.com/urfave/cli"
)

func EnvironmentCommand() cli.Command {
	return cli.Command{
		Name:  "env",
		Usage: "Operations of environment",
		Subcommands: []cli.Command{
			{
				Name:  "clone",
				Usage: "Clone an environment to a new environment",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "source-env",
					},
					cli.StringFlag{
						Name: "target-env",
					},
				},
				Action: CloneEnvironmentAction,
			},
		},
	}
}

func CloneEnvironmentAction(c *cli.Context) error {
	opts := config.EnvUpgradeOpts{
		SourceEnv: c.String("source-env"),
		TargetEnv: c.String("target-env"),
	}

	client, err := rancher.NewClient(cattleUrl, cattleAccessKey, cattleSecret, "")
	if err != nil {
		return err
	}

	return client.CloneProject(opts)
}
