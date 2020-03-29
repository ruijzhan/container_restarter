package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"net"
	"os"
)

type myDockerCli struct {
	*client.Client
}

func (cli *myDockerCli)containerByName(name string) (container *types.Container, ok bool){
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		for _, cName := range container.Names {
			if name == cName[1:] {  //[1:] removed prefix '/' from container's name
				return &container, true
			}
		}
	}
	return nil, false
}

func getDockerClient(host, version string) (*myDockerCli, error)  {

	if host != "unix:///var/run/docker.sock" {
		os.Setenv("DOCKER_HOST", host)
	}
	os.Setenv("DOCKER_API_VERSION", version)

	if cli, err := client.NewClientWithOpts(client.FromEnv); err != nil {
		return nil, err
	} else {
		return &myDockerCli{cli}, nil
	}
}

func ipChangeFunc(domain string) func() bool {
	oldIp := ""
	return func () bool {
		ips, err := net.LookupIP(domain)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot resolve %s\n", domain)
			return  false
		}
		ip := ips[0].String()

		if oldIp == "" {
			oldIp = ip
			return  false
		}

		if ip != oldIp {
			log.Printf("IP address changed from %s to %s", oldIp, ip)
			oldIp = ip
			return  true
		}
		return  false
	}
}

func conditionExec(fExe func() error, fCon func()  bool) error {
	if  con := fCon(); con {
		if err := fExe(); err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func restartContainer(cName string, cli *myDockerCli) func() error {
	return func() error {
		if container, ok := cli.containerByName(cName); ok {
			if err := cli.ContainerRestart(context.Background(), container.ID, nil); err != nil {
				return err
			} else {
				log.Printf("Container %s restarted", cName)
			}
		} else {
			return  fmt.Errorf("container %s not found", cName)
		}
		return nil
	}
}
