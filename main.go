package main

import (
	"github.com/ruijzhan/container_restarter/utils"
)

func main() {
	restarter := utils.NewMyMsgBus()
	restarter.Regist(utils.NewMyContainer(id, container).Restart)

	ipChanged := utils.IPChangeDetector(domainName, interval)
	restarter.NotifiedBy(ipChanged)

	restarter.Run()
}
