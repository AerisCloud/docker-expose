package commands

import (
	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	{
		Name:   "add",
		Usage:  "Add a container to the expose server",
		Action: cmdAdd,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port, p",
				Value: 80,
				Usage: "The port of the container to expose, default to 80",
			},
		},
	},
	{
		Name:   "list",
		Usage:  "List currently running containers and their expose status",
		Action: cmdList,
	},
	{
		Name:   "rm",
		Usage:  "Remove a container from the expose server",
		Action: cmdRm,
	},
}
