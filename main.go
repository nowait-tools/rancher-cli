package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/nowait-tools/rancher-cli/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.3.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Nowait",
			Email: "devops@nowait.com",
		},
	}
	log.SetLevel(log.DebugLevel)
	app.Usage = "The awesome cli you wish existed for Rancher"
	app.Commands = []cli.Command{
		cmd.EnvironmentCommand(),
		cmd.ServiceCommand(),
	}
	err := app.Run(os.Args)

	if err != nil {
		log.Fatalf("command exited with error %v", err)
	}
}
