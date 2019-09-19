package main

import (
	"log"
	"os"
	"strings"

	clilib "github.com/mcastorina/poster/internal/cli"
	"github.com/urfave/cli"
)

func main() {
	cli.SubcommandHelpTemplate =
		strings.Replace(cli.SubcommandHelpTemplate, "command", "resource", -1)

	app := cli.NewApp()
	app.UseShortOptionHandling = true
	app.Commands = cli.Commands{
		{
			Name:  "create",
			Usage: "create a resource",
			Subcommands: []cli.Command{
				{
					Name:   "target",
					Usage:  "create a target resource",
					Action: clilib.CreateTarget,
				},
				{
					Name:   "request",
					Usage:  "create a request resource",
					Action: clilib.CreateRequest,
				},
			},
		},
		{
			Name:   "run",
			Usage:  "run a request",
			Action: clilib.Run,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
