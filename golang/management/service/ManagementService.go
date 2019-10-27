package service

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type ManagementService struct {
	id        uint16
	serviceID *ServiceID
	sm        *ServiceManager
	consoleId *ConsoleId
}

func (srv *ManagementService) Topic() string {
	return "Management Service"
}

func (srv *ManagementService) ID() uint16 {
	return srv.id
}

func (srv *ManagementService) ServiceManager() *ServiceManager {
	return srv.sm
}

func (srv *ManagementService) Init(sm *ServiceManager, id uint16, sid *ServiceID, cid *ConsoleId) {
	srv.sm = sm
	srv.id = id
	srv.consoleId = cid
	srv.serviceID = sid
}

func (srv *ManagementService) ConsoleId() *ConsoleId {
	return srv.consoleId
}

func (srv *ManagementService) ServiceID() *ServiceID {
	return srv.serviceID
}
