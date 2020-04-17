package main

import (
	"testing"
	"time"
)

var (
	ips = []string{
		"1.1.1.1",
		"1.1.1.1",
		"2.2.2.2",
		"2.2.2.2",
		"3.3.3.3",
		"3.3.3.3",
		"3.3.3.3",
		"4.4.4.4",
		"4.4.4.4",
		"5.5.5.5",
		"5.5.5.5",
	}

	changedIPs = []string{
		"2.2.2.2",
		"3.3.3.3",
		"4.4.4.4",
		"5.5.5.5",
	}
)

func fakeLookup() func(_ string) (string, error) {
	var i int
	return func(_ string) (string, error) {
		pi := i
		if pi == len(ips)-1 {
			i = 0
		} else {
			i++
		}
		return ips[pi], nil
	}
}

func TestResolver(t *testing.T) {
	saved := lookup
	defer func() { lookup = saved }()
	lookup = fakeLookup()
	interval = time.Duration(time.Microsecond)

	r := resolver("www.example.com")
	for i := 0; i < len(changedIPs); i++ {
		ip := <-r()
		if changedIPs[i] != ip {
			t.Errorf("IP changed to %s, want %s", ip, changedIPs[i])
		}
	}
}
