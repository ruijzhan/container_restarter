package utils

import (
	"log"
	"net"
	"time"
)

var lookup = func(d string) (string, error) { // define func var for testing
	ips, err := net.LookupIP(d)
	if err != nil {
		return "", err
	}
	return ips[0].String(), nil
}

func Resolver(d string, interval time.Duration) func() chan string {
	oldIP, err := lookup(d)
	if err != nil {
		log.Printf("Warning: %v.", err)
	}
	tick := time.Tick(interval)
	ch := make(chan string)

	return func() chan string {
		go func() {
			for {
				<-tick
				ip, err := lookup(d)

				if err != nil {
					log.Printf("Warning: %v.", err)
				} else {

					if ip != oldIP {
						oldIP = ip
						ch <- ip
						return
					}
				}
			}
		}()

		return ch
	}
}
