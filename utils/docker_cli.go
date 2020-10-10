package utils

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"os"
)

type MyDockerCli struct {
	*client.Client
}

//新建Docker客户端
func NewMyDockerCli(host, version string) (*MyDockerCli, error) {
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}
	if host != "unix:///var/run/docker.sock" {
		os.Setenv("DOCKER_HOST", host)
	}
	if version == "" {
		version = "1.40"
	}
	os.Setenv("DOCKER_API_VERSION", version)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &MyDockerCli{cli}, nil
}

func (cli *MyDockerCli) getContainer(nameOrId string) *types.Container {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range containers {
		for _, cName := range container.Names {
			if nameOrId == cName[1:] { //[1:] removed prefix '/' from container's name
				return &container
			}
		}

		if nameOrId == container.ID {
			return &container
		}
	}
	return nil
}

func (cli *MyDockerCli) restartContainer(nameOrId string) error {

	if container := cli.getContainer(nameOrId); container != nil {
		if err := cli.ContainerRestart(context.Background(), container.ID, nil); err != nil {
			return err
		}
		log.Printf("Container %s restarted", nameOrId)
	} else {
		return fmt.Errorf("container %s not found", nameOrId)
	}
	return nil
}
