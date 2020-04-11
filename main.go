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

	var act func() error
	if *container == "debug" {
		act = func() error { return nil }
	} else {
		act = fRestartC(cli, *container)
	}

	newIP := resolver(*domain)
	for {
		// <-newIP() is blocked till *domain resolved IP changes
		log.Printf("IP address changed to %s", <-newIP())
		act()
	}

}
