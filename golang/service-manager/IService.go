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

type IServiceMessageHandlers interface {
	Init(IService)
	Handler(string) IMessageHandler
	Handlers() []IMessageHandler
}

type IServiceCommands interface {
	Init(IService, IServiceMessageHandlers)
	Commands() map[string]Command
}

type IMessageHandler interface {
	Handle(*Message)
	Topic() string
	Message(*ServiceID, []byte, bool) *Message
	Init()
	Request(interface{}, *ServiceID) interface{}
}
