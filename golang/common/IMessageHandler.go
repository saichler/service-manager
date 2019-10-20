package common

import (
	. "github.com/saichler/messaging/golang/net/protocol"
)

type IMessageHandler interface {
	Send()
	Handle(message *Message)
	Message() *Message
}
