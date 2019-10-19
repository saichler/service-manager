package management

import (
	. "github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/common"
)

type ManagementService struct {
	id uint16
	sm IServiceManager
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

func (srv *ManagementService) Init(sm IServiceManager, id uint16) {
	srv.sm = sm
	srv.id = id
}

func (srv *ManagementService) Handle(message *Message) {

}

func (srv *ManagementService) Start() {

}
