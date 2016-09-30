package main

import (
	"os"

	"github.com/nowait/rancher/cmd"
	"github.com/nowait/rancher/rancher"
	"github.com/urfave/cli"
)

var (
	cattleUrl       = os.Getenv("CATTLE_URL")
	cattleAccessKey = os.Getenv("CATTLE_ACCESS_KEY")
	cattleSecret    = os.Getenv("CATTLE_SECRET_KEY")
)

func main() {
	client, err := rancher.NewClient(cattleUrl, cattleAccessKey, cattleSecret)

	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		cmd.ServiceCommand(client),
	}
	app.Run(os.Args)
}
