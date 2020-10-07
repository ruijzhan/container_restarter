package main

import (
	"flag"
	"os"
	"time"
)

var (
	container  string
	id         string
	domainName string
	interval   time.Duration
	host       string
	version    string
)

func init() {
	flag.StringVar(&container, "c", "", "Name of the container to restart")
	flag.StringVar(&id, "id", "", "ID of the container to restart")
	flag.StringVar(&domainName, "d", "", "Domain name to watch IP change")
	flag.DurationVar(&interval, "t", time.Duration(10*time.Second),
		"Time interval to check IP change on domain name")
	flag.StringVar(&host, "h", "unix:///var/run/docker.sock",
		"Docker server host")
	flag.StringVar(&version, "v", "1.40", "Docker API version")

	flag.Parse()
	validate()
}

func validate() {
	if (container == "" && id == "") || domainName == "" {
		flag.Usage()
		os.Exit(1)
	}
}
