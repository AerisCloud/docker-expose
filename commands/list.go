package commands

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"github.com/ttacon/chalk"

	"github.com/aeriscloud/docker-expose/expose"
)

// find the shortest name for a container if it has several
func bestName(names []string) string {
	var best string
	for _, name := range names {
		if best == "" || len(name) < len(best) {
			best = name
		}
	}
	return best[1:]
}

func cmdList(c *cli.Context) {
	exposeClient, err := expose.NewExposeFromConf()
	if err != nil {
		log.Fatalln(err)
	}

	exposedHosts, err := exposeClient.List(true)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalln(err)
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%-20s%-48s%-64s%s\n", "CONTAINER ID", "CONTAINER NAME", "URL", "STATUS")
	for _, container := range containers {
		name := bestName(container.Names)
		exposedHost, isExposed := exposedHosts.Find(name)
		url := exposeClient.UserURL(name)
		status := chalk.Blue.String() + "NOT EXPOSED" + chalk.Reset.String()
		if isExposed {
			if exposedHost.Status == 0 {
				status = chalk.Red.String() + "DOWN" + chalk.Reset.String()
			}

			if exposedHost.Status == 1 {
				status = chalk.Yellow.String() + "MAINTENANCE" + chalk.Reset.String()
			}

			if exposedHost.Status == 2 {
				status = chalk.Green.String() + "AVAILABLE" + chalk.Reset.String()
			}
		}

		fmt.Printf("%-20s%-48s%-64s%s\n", container.ID[:12], name, url, status)
	}
}
