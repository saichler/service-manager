package service

import (
	. "github.com/saichler/messaging/golang/net/protocol"
)

type Service interface {
	Topic() string
	ID() uint16
	ServiceManager() *ServiceManager
	Init(*ServiceManager, uint16)
	Handle(message *Message)
	Start()
}
