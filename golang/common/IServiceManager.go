package common

import (
	. "github.com/saichler/console/golang/console"
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/messaging/golang/net/protocol"
)

type IServiceManager interface {
	Shutdown()
	ConsoleId() *ConsoleId
	Console() *Console
	Services() []IService
	NetworkID() *NetworkID
	Publish(string, IService, []byte) error
	Send(string, IService, *ServiceID, []byte) error
	Uplink(ip string)
	ScheduleMessage(IMessageHandler, int64, int64)
	NewMessage(string, IService, *ServiceID, []byte) *Message
}
