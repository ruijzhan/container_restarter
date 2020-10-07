package main

import (
	"github.com/ruijzhan/container_restarter/utils"
	"log"
)

func main() {
	restarter := utils.NewRestarter()
	err := restarter.Regist(utils.NewMyContainer(id, container).Restart)
	if err != nil {
		log.Fatal(err)
	}

	ipChanged := utils.IPChangeDetector(domainName, interval)
	restarter.NotifiedBy(ipChanged)

	if err := restarter.Run(); err != nil {
		log.Fatal(err)
	}
}
