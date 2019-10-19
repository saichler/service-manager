package common

import . "github.com/saichler/messaging/golang/net/protocol"

type IService interface {
	Topic() string
	ID() uint16
	ServiceManager() IServiceManager
	Init(IServiceManager, uint16)
	Handle(message *Message)
	Start()
}
