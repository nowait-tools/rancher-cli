package main

import (
	"os"

	"github.com/nowait/rancher-cli/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		cmd.ServiceCommand(),
	}
	app.Run(os.Args)
}
