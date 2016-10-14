package main

import (
	"log"
	"os"

	"github.com/nowait/rancher-cli/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.3.0-rc2"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Nowait",
			Email: "devops@nowait.com",
		},
	}
	app.Usage = "The awesome cli you wish existed for Rancher"
	app.Commands = []cli.Command{
		cmd.ServiceCommand(),
	}
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
