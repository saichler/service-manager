package service

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/common"
	"github.com/saichler/service-manager/golang/management/commands"
	"github.com/saichler/service-manager/golang/management/model"
)

type ManagementService struct {
	id             uint16
	serviceID      *ServiceID
	sm             IServiceManager
	consoleId      *ConsoleId
	serviceNetwork *model.ServiceNetwork
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

func (srv *ManagementService) Init(sm IServiceManager, id uint16, sid *ServiceID, cid *ConsoleId) []IMessageHandler {
	srv.serviceNetwork = model.NewServiceNetwork()
	srv.sm = sm
	srv.id = id
	srv.consoleId = cid
	srv.serviceID = sid
	sm.Console().RegisterCommand(commands.NewShutdown(srv))
	sm.Console().RegisterCommand(commands.NewUplink(srv))

	handlers := make([]IMessageHandler, 0)
	ping := NewPingMH(srv)
	handlers = append(handlers, ping)
	sm.ScheduleMessage(ping, 5, 0)
	return handlers
}

func (srv *ManagementService) ConsoleId() *ConsoleId {
	return srv.consoleId
}

func (srv *ManagementService) ServiceID() *ServiceID {
	return srv.serviceID
}
