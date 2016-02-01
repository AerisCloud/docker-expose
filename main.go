package main

import (
	"os"
	"path"

	"github.com/aeriscloud/docker-expose/commands"

	"github.com/codegangsta/cli"
)

var (
	VERSION = "0.1.0"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = "Wizcorp Inc"
	app.Usage = "Expose docker containers on aeris.cd"
	app.Commands = commands.Commands
	app.Version = VERSION
	app.Run(os.Args)
}
