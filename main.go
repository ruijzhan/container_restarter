package main

import (
	"flag"
	"github.com/ruijzhan/container_restarter/utils"
	"log"
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
}

func main() {
	//命令行参数初始化
	flag.Parse()

	//参数判断: container和id 二选一, domainName 必填
	if (container == "" && id == "") || domainName == "" {
		flag.Usage()
		os.Exit(1)
	}

	//Type: channel
	//用于存取解析的IP
	ipChanged := utils.Resolver(domainName, interval)
	restarter := utils.NewMyMsgBus()
	restarter.Regist(utils.NewMyContainer(id, container))

	for {
		// <-ipChanged() is blocked till *domainName resolved IP changes
		select {
		case newIP := <-ipChanged():
			log.Printf("IP address changed to %s", newIP)
			restarter.Notify()
		}
	}

}
