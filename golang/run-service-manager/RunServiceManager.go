package main

import (
	"github.com/saichler/service-manager/golang/management/commands"
	message_handlers "github.com/saichler/service-manager/golang/management/message-handlers"
	"github.com/saichler/service-manager/golang/management/service"
	"github.com/saichler/service-manager/golang/service-manager"
)

func main() {
	serviceManager, e := service_manager.NewServiceManager()
	if e != nil {
		return
	}

	serviceManager.Console().RegisterCommand(commands.NewLS(serviceManager))
	serviceManager.Console().RegisterCommand(commands.NewCD(serviceManager))
	serviceManager.AddService(NewService())
	serviceManager.LoadService("/home/saichler/datasand/src/github.com/saichler/service-manager/golang/file-manager/plugin/FileService.so")
	serviceManager.WaitForShutdown()
}

func NewService() (service_manager.IService, service_manager.IServiceCommands, service_manager.IServiceMessageHandlers) {
	s := &service.ManagementService{}
	h := &message_handlers.ManagementHandlers{}
	h.Init(s)
	c := &commands.ManagementCommands{}
	c.Init(s, h)
	return s, c, h
}
