package cmd

import (
	"os"

	"github.com/nowait/rancher-cli/rancher"
	"github.com/urfave/cli"
)

var (
	cattleUrl       = os.Getenv("CATTLE_URL")
	cattleAccessKey = os.Getenv("CATTLE_ACCESS_KEY")
	cattleSecret    = os.Getenv("CATTLE_SECRET_KEY")
)

func ServiceCommand() cli.Command {
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
					client, err := rancher.NewClient(cattleUrl, cattleAccessKey, cattleSecret)
					if err != nil {
						return err
					}
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
						Name: "service-like",
					},
					cli.StringFlag{
						Name: "tag",
					},
				},
				Action: UpgradeCodeAction,
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
					client, err := rancher.NewClient(cattleUrl, cattleAccessKey, cattleSecret)

					if err != nil {
						return err
					}
					_, err = client.FinishServiceUpgrade(c.String("service"))
					return err
				},
			},
		},
	}
}

func UpgradeCodeAction(c *cli.Context) error {
	client, err := rancher.NewClient(cattleUrl, cattleAccessKey, cattleSecret)
	if err != nil {
		return err
	}

	if name := c.String("service-like"); name != "" {
		return client.UpgradeServiceWithNameLike(name, c.String("tag"))
	}
	return client.UpgradeServiceCodeVersion(c.String("service"), c.String("tag"))
}
