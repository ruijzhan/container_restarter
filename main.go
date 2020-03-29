package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var (
	container = flag.String("c", "", "Name of container to restart")
	domain    = flag.String("d", "", "Domain name to watch IP change")
	interval  = flag.Duration("i", time.Duration(10*time.Second), "Time interval to check IP change on domain")
	host      = flag.String("h", "unix:///var/run/docker.sock", "docker server host")
	version   = flag.String("v", "1.40", "Docker API version")
)

func main() {
	flag.Parse()
	if *container == "" || *domain == "" {
		flag.Usage()
		os.Exit(1)
	}

	cli, err := myDockerClient(*host, *version)
	if err != nil {
		log.Fatal(err)
	}

	restartCondition := ipChanged(*domain)
	restartCondition() //init resolve result
	restartContainer := run2RestartContainer(cli, *container)
	for {
		runIfTrue(restartContainer, restartCondition)
		time.Sleep(*interval)
	}

}
