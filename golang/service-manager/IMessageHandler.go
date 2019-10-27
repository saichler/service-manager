package service_manager

import (
	. "github.com/saichler/messaging/golang/net/protocol"
)

type IMessageHandler interface {
	Send()
	Handle(message *Message)
	Topic() string
	Message() *Message
	Init()
}
