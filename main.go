package main

import (
	"flag"
	"github.com/ruijzhan/container_restarter/tools"
	mbus "github.com/vardius/message-bus"
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

type myContainer struct {
	dockerCli *tools.MyDockerCli
	id        string
	name      string
}

func NewMyContainer(id, name string) *myContainer {
	cli, err := tools.MyDockerClient(host, version)
	if err != nil {
		log.Fatal(err)
	}

	return &myContainer{
		dockerCli: cli,
		id:        id,
		name:      name,
	}
}

func (c *myContainer) restart() {
	c.dockerCli.RestartContainer(c.id, c.name)
}

type myMsgBus struct {
	mbus.MessageBus
	topic string
}

func NewMyMsgBus() *myMsgBus {
	return &myMsgBus{
		mbus.New(10),
		"ipChanged",
	}
}

func (m *myMsgBus) regist(c *myContainer) {
	m.Subscribe(m.topic, c.restart)
}

func (m *myMsgBus) notify() {
	m.Publish(m.topic)
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
	newIP := tools.Resolver(domainName, interval)
	bus := NewMyMsgBus()
	bus.regist(NewMyContainer(id, container))

	for {
		// <-newIP() is blocked till *domainName resolved IP changes
		select {
		case ip := <-newIP():
			log.Printf("IP address changed to %s", ip)
			bus.notify()
		}
	}

}
