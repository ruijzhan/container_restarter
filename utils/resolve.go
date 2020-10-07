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

func IPChangeDetector(d string, interval time.Duration) <-chan string {
	oldIP, err := lookup(d)
	if err != nil {
		log.Printf("Warning: %v.", err)
	}
	tick := time.Tick(interval)
	ch := make(chan string)

	go func() {
		for {
			select {
			case <-tick:
				ip, err := lookup(d)
				if err != nil {
					log.Printf("Warning: %v.", err)
					continue
				}
				if ip != oldIP {
					oldIP = ip
					ch <- ip
				}
			}
		}
	}()

	return ch
}
