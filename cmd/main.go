package main

import (
	"log"
	"os"
	"strings"

	clilib "github.com/mcastorina/poster/internal/cli"
	"github.com/mcastorina/poster/internal/store"
	"github.com/urfave/cli"
)

func main() {
	cli.SubcommandHelpTemplate =
		strings.Replace(cli.SubcommandHelpTemplate, "command", "resource", -1)
	cli.SubcommandHelpTemplate =
		strings.Replace(cli.SubcommandHelpTemplate, "COMMANDS", "RESOURCES", -1)

	app := cli.NewApp()
	app.UseShortOptionHandling = true
	app.Commands = cli.Commands{
		{
			Name:  "create",
			Usage: "create a resource",
			Subcommands: []cli.Command{
				{
					Name:    "target",
					Aliases: []string{"t"},
					Usage:   "create a target resource",
					Action:  clilib.CreateTarget,
				},
				{
					Name:    "request",
					Aliases: []string{"r", "req"},
					Usage:   "create a request resource",
					Action:  clilib.CreateRequest,
				},
			},
		},
		{
			Name:  "get",
			Usage: "view resources",
			Subcommands: []cli.Command{
				{
					Name:    "target",
					Aliases: []string{"t", "targets"},
					Usage:   "view target resources",
					Action:  clilib.GetTarget,
				},
				{
					Name:    "request",
					Aliases: []string{"r", "req", "requests", "reqs"},
					Usage:   "view request resources",
					Action:  clilib.GetRequest,
				},
			},
		},
		{
			Name:   "run",
			Usage:  "run a request",
			Action: clilib.Run,
		},
	}

	store.InitDB()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
