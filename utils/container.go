package utils

import (
	"log"
)

type MyContainer struct {
	dockerCli *MyDockerCli
	id        string
	name      string
}

func NewMyContainer(id, name string) *MyContainer {
	cli, err := NewMyDockerCli("", "")
	if err != nil {
		log.Fatal(err)
	}

	return &MyContainer{
		dockerCli: cli,
		id:        id,
		name:      name,
	}
}

func (c *MyContainer) restart() {
	c.dockerCli.restartContainer(c.id, c.name)
}
