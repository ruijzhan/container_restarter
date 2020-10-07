package utils

import (
	mbus "github.com/vardius/message-bus"
	"log"
)

const TOPIC = "ipChanged"

type Restarter interface {
	Regist(func())
	NotifiedBy(<-chan string)
	Run()
}

type myMsgBus struct {
	mbus.MessageBus
	notifier <-chan string
}

func NewMyMsgBus() Restarter {
	return &myMsgBus{
		mbus.New(10),
		nil,
	}
}

func (m *myMsgBus) Regist(f func()) {
	m.Subscribe(TOPIC, f)
}

func (m *myMsgBus) NotifiedBy(ch <-chan string) {
	m.notifier = ch
}

func (m *myMsgBus) Run() {
	if m.notifier == nil {
		log.Fatal("Notifier not set")
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
