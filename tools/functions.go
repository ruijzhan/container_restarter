package tools

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"net"
	"os"
	"time"
)

type myDockerCli struct {
	*client.Client
}

//新建Docker客户端
func MyDockerClient(host, version string) (*myDockerCli, error) {
	if host != "unix:///var/run/docker.sock" {
		os.Setenv("DOCKER_HOST", host)
	}
	os.Setenv("DOCKER_API_VERSION", version)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &myDockerCli{cli}, nil
}

var lookup = func(d string) (string, error) { // define func var for testing
	ips, err := net.LookupIP(d)
	if err != nil {
		return "", err
	}
	return ips[0].String(), nil
}

func Resolver(d string, interval time.Duration) func() chan string {
	oldIP, err := lookup(d)
	if err != nil {
		log.Printf("Warning: %v.", err)
	}
	tick := time.Tick(interval)
	ch := make(chan string)

	return func() chan string {
		go func() {
			for {
				<-tick
				ip, err := lookup(d)

				if err != nil {
					log.Printf("Warning: %v.", err)
				} else {

					if ip != oldIP {
						oldIP = ip
						ch <- ip
						return
					}
				}
			}
		}()

		return ch
	}
}

//判断给定容器名称对应的目标容器是否存在，存在则返回对应的容器对象 &container
func (cli *myDockerCli) getContainerByName(name string) (container *types.
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
func (cli *myDockerCli) getContainerByID(id string) (container *types.
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
func RestartContainerByName(cli *myDockerCli, name string) error {
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
func RestartContainerByID(cli *myDockerCli, id string) error {
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

func RestartContainer(cli *myDockerCli, id, name string) error {
	if id != "" {
		return RestartContainerByID(cli, id)
	}
	if name != "" {
		return RestartContainerByName(cli, name)
	}
	return fmt.Errorf("contaienr name or id no set")
}
