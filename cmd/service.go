package cmd

import (
	"github.com/nowait/rancher/rancher"
	"github.com/urfave/cli"
)

func ServiceCommand(client *rancher.Client) cli.Command {
	return cli.Command{
		Name:  "service",
		Usage: "Operations on services",
		Subcommands: []cli.Command{
			{
				Name:  "upgrade-runtime",
				Usage: "Upgrade the runtime tag of the service",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "service",
					},
					cli.StringFlag{
						Name: "tag",
					},
				},
				Action: func(c *cli.Context) error {
					return client.UpgradeServiceVersion(c.String("service"), c.String("tag"))
				},
			},
			{
				Name:  "upgrade-code",
				Usage: "Upgrade the code tag of the service",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "service",
					},
					cli.StringFlag{
						Name: "tag",
					},
				},
				Action: func(c *cli.Context) error {
					return client.UpgradeServiceCodeVersion(c.String("service"), c.String("tag"))
				},
			},
			{
				Name:  "upgrade-finish",
				Usage: "",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "service",
					},
				},
				Action: func(c *cli.Context) error {
					_, err := client.FinishServiceUpgrade(c.String("service"))
					return err
				},
			},
		},
	}
}
