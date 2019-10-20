package service

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/common"
	"github.com/saichler/service-manager/golang/management/commands"
)

type ManagementService struct {
	id        uint16
	serviceID *ServiceID
	sm        IServiceManager
	consoleId *ConsoleId
}

func (srv *ManagementService) Topic() string {
	return "Management Service"
}

func (srv *ManagementService) ID() uint16 {
	return srv.id
}

func (srv *ManagementService) ServiceManager() IServiceManager {
	return srv.sm
}

func (srv *ManagementService) Init(sm IServiceManager, id uint16, consoleId *ConsoleId) {
	srv.sm = sm
	srv.id = id
	srv.consoleId = consoleId
	srv.serviceID = NewServiceID(sm.NetworkID(), srv.Topic(), id)
	sm.Console().RegisterCommand(commands.NewShutdown(srv))
	sm.Console().RegisterCommand(commands.NewCD(srv))
	sm.Console().RegisterCommand(commands.NewUplink(srv))
	msg := NewInventoryMessage(srv)
	sm.ScheduleMessage(msg, 5, 0)
}

func (srv *ManagementService) Handle(message *Message) {

}

func (srv *ManagementService) Start() {

}

func (srv *ManagementService) ConsoleId() *ConsoleId {
	return srv.consoleId
}

func (srv *ManagementService) ServiceID() *ServiceID {
	return srv.serviceID
}
