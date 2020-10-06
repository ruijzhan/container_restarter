package utils

import mbus "github.com/vardius/message-bus"

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

func (m *myMsgBus) Regist(c *MyContainer) {
	m.Subscribe(m.topic, c.restart)
}

func (m *myMsgBus) Notify() {
	m.Publish(m.topic)
}
