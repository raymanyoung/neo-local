package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/CityOfZion/neo-local/cli/logger"
	"github.com/CityOfZion/neo-local/cli/stack"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// CheckDockerRunning pings the Docker client to check that it is running.
func CheckDockerRunning(ctx context.Context, cli *client.Client) bool {
	response, err := cli.Ping(ctx)
	if err != nil {
		return false
	}

	return response.APIVersion != ""
}

// PullDockerImages pulls each Docker image required for the stack
// sequentially.
func PullDockerImages(ctx context.Context, cli *client.Client) error {
	log.Println("Pulling Docker images")

	services, err := stack.Services()
	if err != nil {
		return err
	}

	for _, service := range services {
		prefix := fmt.Sprintf("↪  %s", service.ImageName())
		s := logger.NewSpinner(prefix)
		s.Start()

		reader, err := cli.ImagePull(
			ctx, service.ImageName(), types.ImagePullOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"Error when pulling Docker image '%s': %s",
				service.ImageName(),
				err,
			)
		}

		var b bytes.Buffer
		_, err = io.Copy(&b, reader)
		if err != nil {
			return fmt.Errorf(
				"Error when pulling Docker image '%s': %s",
				service.ImageName(),
				err,
			)
		}

		defer reader.Close()
		s.Stop()
	}

	return nil
}

// FetchContainerReferences finds the container ID for each service within the
// stack.
func FetchContainerReferences(ctx context.Context, cli *client.Client) (map[string]*types.Container, error) {
	containerReferences := map[string]*types.Container{}

	serviceContainerNames, err := stack.ServiceContainerNames()
	if err != nil {
		return nil, err
	}

	for _, serviceContainerName := range serviceContainerNames {
		containerReferences[serviceContainerName] = nil
	}

	containers, err := cli.ContainerList(
		ctx,
		types.ContainerListOptions{
			All: true,
		},
	)
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		containerName, ok := isNeoLocalContainer(container.Names)
		if !ok {
			continue
		}

		for _, serviceContainerName := range serviceContainerNames {
			if containerName == serviceContainerName {
				containerReferences[serviceContainerName] = &container
				break
			}
		}
	}

	return containerReferences, nil
}

func isNeoLocalContainer(containerNames []string) (string, bool) {
	for _, containerName := range containerNames {
		containerName = strings.Replace(containerName, "/", "", -1)
		if strings.HasPrefix(containerName, stack.ContainerNamePrefix) {
			return containerName, true
		}
	}

	return "", false
}
