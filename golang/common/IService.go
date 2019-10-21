package common

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/messaging/golang/net/protocol"
)

type IService interface {
	Topic() string
	ID() uint16
	ConsoleId() *ConsoleId
	ServiceManager() IServiceManager
	Init(IServiceManager, uint16, *ServiceID, *ConsoleId) []IMessageHandler
	ServiceID() *ServiceID
}
