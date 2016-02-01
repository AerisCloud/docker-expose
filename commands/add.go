package commands

import (
	"fmt"
	"log"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"

	"github.com/aeriscloud/docker-expose/expose"
)

func cmdAdd(c *cli.Context) {
	args := c.Args()

	if len(args) != 1 {
		cli.ShowAppHelp(c)
		log.Fatalln("missing argument: container name")
	}

	exposeClient, err := expose.NewExposeFromConf()
	if err != nil {
		log.Fatalln(err)
	}

	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalln(err)
	}

	container, err := client.InspectContainer(args[0])
	if err != nil {
		log.Fatalln(err)
	}

	port := c.Int("port")
	dockerPort := docker.Port(fmt.Sprintf("%d/tcp", port))
	portBindings := container.NetworkSettings.Ports[dockerPort]

	if len(portBindings) == 0 {
		log.Fatalf("port %d is not exposed by the container\n", port)
	}

	hostPort, _ := strconv.ParseInt(portBindings[0].HostPort, 10, 32)
	err = exposeClient.Add(container.Name[1:], int(hostPort))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Container %s successfully exposed!\n", container.Name[1:])
	fmt.Println(exposeClient.UserURL(container.Name[1:]))
}
