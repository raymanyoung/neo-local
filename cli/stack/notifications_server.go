package stack

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/CityOfZion/neo-local/cli/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	packr "github.com/gobuffalo/packr/v2"
)

func initFile(file string, box packr.Box) error {
	dirPath, err := config.DirPath()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%s", dirPath, file)

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		fileContent, err := box.Find(file)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, fileContent, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewNotificationsServer creates a new service for the cityofzion/neo-python image.
func NewNotificationsServer() (Service, error) {
	box := packr.New("notificationsBox", "./../../notifications-server")
	initFile("notifications-server.config.json", *box)

	dirPath, err := config.DirPath()
	if err != nil {
		return Service{}, err
	}

	binds := []string{
		fmt.Sprintf("%s:/configs", dirPath),
	}

	return Service{
		Author: "cityofzion",
		ContainerConfig: &container.Config{
			Cmd: []string{
				"/bin/sh",
				"-c",
				"/bin/cp /configs/notifications-server.config.json /neo-python/custom-config.json && /usr/bin/python3 /neo-python/neo/bin/api_server.py --config /neo-python/custom-config.json --port-rest 8080",
			},
			Env: []string{
				"NOTIFICATIONS_SERVER=notifications-server",
			},
			ExposedPorts: map[nat.Port]struct{}{
				"8080/tcp": {},
			},
		},
		HostConfig: &container.HostConfig{
			Binds: binds,
			PortBindings: map[nat.Port][]nat.PortBinding{
				"8080/tcp": {
					{
						HostIP:   "0.0.0.0",
						HostPort: "8080",
					},
				},
			},
			Privileged:      false,
			PublishAllPorts: false,
		},
		Image: "neo-python",
		Name:  "notifications-server",
		Tag:   "v0.8.2",
	}, nil
}
