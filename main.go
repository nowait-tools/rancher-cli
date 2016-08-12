package main

import (
	"os"

	"github.com/nowait/rancher/rancher"
	"github.com/urfave/cli"
)

const (
	cattleUrl       = "http://192.168.99.101:8080"
	cattleAccessKey = "F7D521403075A7D29088"
	cattleSecret    = "en33gY5NgGiyRBu5LrdqjGuCTEuxCLsfpk2f1Ndr"
)

func main() {
	client, err := rancher.NewClient(cattleUrl, cattleAccessKey, cattleSecret)

	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "upgrade",
			Usage: "",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "service"},
				cli.StringFlag{Name: "runtime-tag"},
			},
			Action: func(c *cli.Context) error {
				return client.UpgradeServiceVersion(c.String("service"), c.String("runtime-tag"))
			},
		},
		{
			Name:  "upgrade-code",
			Usage: "",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "service"},
				cli.StringFlag{Name: "code-tag"},
			},
			Action: func(c *cli.Context) error {
				return client.UpgradeServiceCodeVersion(c.String("service"), c.String("code-tag"))
			},
		},
		{
			Name:  "upgrade-finish",
			Usage: "",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "service"},
			},
			Action: func(c *cli.Context) error {
				_, err := client.FinishServiceUpgrade(c.String("service"))
				return err
			},
		},
	}
	app.Run(os.Args)
}
