package cmd

import (
	"os"
	"time"

	"github.com/nowait/rancher-cli/rancher"
	"github.com/urfave/cli"
)

var (
	cattleUrl       = os.Getenv("CATTLE_URL")
	cattleAccessKey = os.Getenv("CATTLE_ACCESS_KEY")
	cattleSecret    = os.Getenv("CATTLE_SECRET_KEY")

	defaultUpgradeInterval = 10 * time.Second
)

func ServiceCommand() cli.Command {
	return cli.Command{
		Name:  "service",
		Usage: "Operations on services",
		Subcommands: []cli.Command{
			// {
			// 	Name:  "upgrade-runtime",
			// 	Usage: "Upgrade the runtime tag of the service",
			// 	Flags: []cli.Flag{
			// 		cli.StringFlag{
			// 			Name: "service",
			// 		},
			// 		cli.StringFlag{
			// 			Name: "tag",
			// 		},
			// 	},
			// 	Action: func(c *cli.Context) error {
			// 		client, err := rancher.NewClient(cattleUrl, cattleAccessKey, cattleSecret)
			// 		if err != nil {
			// 			return err
			// 		}
			// 		opts := rancher.UpgradeOpts{
			// 			Service:    c.String("service"),
			// 			RuntimeTag: c.String("tag"),
			// 		}
			// 		return client.UpgradeServiceVersion(opts)
			// 	},
			// },
			{
				Name:  "upgrade",
				Usage: "Upgrade a service",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "service",
					},
					cli.StringFlag{
						Name: "service-like",
					},
					cli.StringFlag{
						Name: "runtime-tag",
					},
					cli.StringFlag{
						Name: "code-tag",
					},
					cli.Int64Flag{
						Name:  "interval",
						Usage: "Interval between starting new containers and stopping old ones",
					},
					cli.BoolFlag{
						Name:  "wait",
						Usage: "Wait for the upgrade to fully complete",
					},
				},
				Action: UpgradeAction,
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

func UpgradeAction(c *cli.Context) error {
	client, err := rancher.NewClient(cattleUrl, cattleAccessKey, cattleSecret)
	if err != nil {
		return err
	}
	interval := time.Duration(0)
	if interval = time.Duration(c.Int64("interval")); interval == 0 {
		interval = defaultUpgradeInterval
	}

	opts := rancher.UpgradeOpts{
		Interval:    interval,
		ServiceLike: c.String("service-like"),
		Service:     c.String("service"),
		CodeTag:     c.String("code-tag"),
		RuntimeTag:  c.String("runtime-tag"),
		Wait:        c.Bool("wait"),
	}
	if name := opts.ServiceLike; name != "" {
		return client.UpgradeServiceWithNameLike(opts)
	}
	return client.UpgradeService(opts)
}
