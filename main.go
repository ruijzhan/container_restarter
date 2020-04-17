package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var (
	container string
	domain    string
	interval  time.Duration
	host      string
	version   string
)

func init() {
	flag.StringVar(&container, "c", "", "Name of container to restart")
	flag.StringVar(&domain, "d", "", "Domain name to watch IP change")
	flag.DurationVar(&interval, "i", time.Duration(10*time.Second), "Time interval to check IP change on domain")
	flag.StringVar(&host, "h", "unix:///var/run/docker.sock", "docker server host")
	flag.StringVar(&version, "v", "1.40", "Docker API version")
}

func main() {
	flag.Parse()
	if container == "" || domain == "" {
		flag.Usage()
		os.Exit(1)
	}

	cli, err := myDockerClient(host, version)
	if err != nil {
		log.Fatal(err)
	}

	var act func() error
	if container == "debug" {
		act = func() error { return nil }
	} else {
		act = fRestartC(cli, container)
	}

	newIP := resolver(domain)
	for {
		// <-newIP() is blocked till *domain resolved IP changes
		log.Printf("IP address changed to %s", <-newIP())
		act()
	}

}
