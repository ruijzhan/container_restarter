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

//判断给定容器名称对应的目标容器是否存在，存在则返回对应的容器对象 &container
func (cli *MyDockerCli) getContainerByName(name string) (container *types.
	Container, ok bool) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		//根据名称判断目标容器是否存在
		for _, cName := range container.Names {
			if name == cName[1:] { //[1:] removed prefix '/' from container's name
				return &container, true
			}
		}
	}
	return nil, false
}

//判断给定容器id对应的目标容器是否存在，存在则返回对应的容器对象 &container
func (cli *MyDockerCli) getContainerByID(id string) (container *types.
	Container, ok bool) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		//根据id判断目标容器是否存在
		if id == container.ID {
			return &container, true
		}
	}
	return nil, false
}

//根据名称重启container
func (cli *MyDockerCli) restartContainerByName(name string) error {
	if container, ok := cli.getContainerByName(name); ok {
		//重启目标容器
		if err := cli.ContainerRestart(context.Background(), container.ID, nil); err != nil {
			return err
		}
		log.Printf("Container %s restarted", name)
	} else {
		return fmt.Errorf("container %s not found", name)
	}
	return nil
}

//根据id重启container
func (cli *MyDockerCli) restartContainerByID(id string) error {
	if container, ok := cli.getContainerByID(id); ok {
		//重启目标容器
		if err := cli.ContainerRestart(context.Background(), container.ID, nil); err != nil {
			return err
		}
		log.Printf("Container ID %s restarted", id)
	} else {
		return fmt.Errorf("container ID %s not found", id)
	}
	return nil
}

func (cli *MyDockerCli) restartContainer(id, name string) error {
	if id != "" {
		return cli.restartContainerByID(id)
	}
	if name != "" {
		return cli.restartContainerByName(name)
	}
	return fmt.Errorf("contaienr name or id no set")
}
