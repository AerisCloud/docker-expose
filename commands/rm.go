package commands

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"

	"github.com/aeriscloud/docker-expose/expose"
)

func cmdRm(c *cli.Context) {
	args := c.Args()

	if len(args) != 1 {
		cli.ShowAppHelp(c)
		log.Fatalln("missing argument: container name")
	}

	exposeClient, err := expose.NewExposeFromConf()
	if err != nil {
		log.Fatalln(err)
	}

	exposedHosts, err := exposeClient.List(true)
	if err != nil {
		log.Fatalln(err)
	}

	name := args[0]
	if _, isExposed := exposedHosts.Find(name); !isExposed {
		log.Fatalf("container %s is not exposed\n", name)
	}

	err = exposeClient.Delete(name)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Container %s successfully removed!\n", name)
}
