package utils

import (
	"fmt"
	mbus "github.com/vardius/message-bus"
	"log"
)

const TOPIC = "ipChanged"

type Restarter interface {
	Regist(func()) error
	NotifiedBy(<-chan string)
	Run() error
}

type myMsgBus struct {
	mbus.MessageBus
	notifier <-chan string
}

func NewRestarter() Restarter {
	return &myMsgBus{
		mbus.New(10),
		nil,
	}
}

func (m *myMsgBus) Regist(f func()) error {
	return m.Subscribe(TOPIC, f)
}

func (m *myMsgBus) NotifiedBy(ch <-chan string) {
	m.notifier = ch
}

func (m *myMsgBus) Run() error {
	if m.notifier == nil {
		return fmt.Errorf("notifier not set")
	}
	for {
		// <-ipChanged() is blocked till *domainName resolved IP changes
		select {
		case newIP := <-m.notifier:
			log.Printf("IP address changed to %s", newIP)
			m.Publish(TOPIC)
		}
	}
}
