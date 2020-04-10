package main

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

func myDockerClient(host, version string) (*myDockerCli, error) {
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

func detectIPChange(d string, i time.Duration) {
	oldIP := ""
	for {
		ips, err := net.LookupIP(d)
		if err != nil {
			log.Printf("Warning: %v.", err)
		}
		ip := ips[0].String()

		if oldIP == "" {
			oldIP = ip
		}

		if ip != oldIP {
			//log.Printf("IP address changed from %s to %s", oldIP, ip)
			oldIP = ip
			changedIP <- ip
		}
		time.Sleep(i)
	}

}

func (cli *myDockerCli) getContainer(name string) (container *types.Container, ok bool) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		for _, cName := range container.Names {
			if name == cName[1:] { //[1:] removed prefix '/' from container's name
				return &container, true
			}
		}
	}
	return nil, false
}

func run2RestartContainer(cli *myDockerCli, name string) func() error {
	return func() error {
		if container, ok := cli.getContainer(name); ok {
			if err := cli.ContainerRestart(context.Background(), container.ID, nil); err != nil {
				return err
			}
			log.Printf("Container %s restarted", name)
		} else {
			return fmt.Errorf("container %s not found", name)
		}
		return nil
	}
}
