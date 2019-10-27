package service_manager

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/messaging/golang/net/protocol"
)

type IService interface {
	Topic() string
	ID() uint16
	ConsoleId() *ConsoleId
	ServiceManager() *ServiceManager
	Init(*ServiceManager, uint16, *ServiceID, *ConsoleId)
	ServiceID() *ServiceID
}

type IServiceCommands interface {
	Commands(IService) []Command
}

type IServiceMessageHandlers interface {
	Handlers(IService) []IMessageHandler
}
