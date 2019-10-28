package service_manager

import (
	. "github.com/saichler/messaging/golang/net/protocol"
)

type IMessageHandler interface {
	Handle(*Message)
	Topic() string
	Message(*ServiceID, []byte, bool) *Message
	Init()
}
