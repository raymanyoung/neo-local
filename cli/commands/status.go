package commands

import (
	"errors"
	"log"

	"github.com/CityOfZion/neo-local/cli/services"
	"github.com/urfave/cli"

	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type (
	// Status is the CLI command for checking the status of containers in the
	// development environment.
	Status struct{}
)

// NewStatus creates a new Status.
func NewStatus() Status {
	return Status{}
}

// ToCommand generates the CLI command struct.
func (s Status) ToCommand() cli.Command {
	return cli.Command{
		Action:  s.action(),
		Aliases: []string{"ps"},
		Flags:   s.flags(),
		Name:    "status",
		Usage:   "Output overall health of network",
	}
}

func (s Status) action() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			return errors.New("Unable to create Docker client")
		}

		ok := services.CheckDockerRunning(ctx, cli)
		if !ok {
			return errors.New("Docker is not running")
		}

		containerReferences, err := services.FetchContainerReferences(ctx, cli)
		if err != nil {
			return err
		}

		for serviceContainerName, container := range containerReferences {
			if container == nil {
				log.Printf("'%s' does not exist", serviceContainerName)
			} else {
				log.Printf(
					"'%s' in '%s' state (#%s)",
					serviceContainerName,
					container.State,
					container.ID[:10],
				)
			}
		}

		return nil
	}
}

func (s Status) flags() []cli.Flag {
	return []cli.Flag{}
}
