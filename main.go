package main

import (
	"log"
	"os"

	"github.com/nowait/rancher-cli/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		cmd.ServiceCommand(),
	}
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
