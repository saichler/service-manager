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
	serviceManager.AddService(&service.ManagementService{}, &commands.ManagementCommands{}, &message_handlers.ManagementHandlers{})
	serviceManager.LoadService("/home/saichler/datasand/src/github.com/saichler/service-manager/golang/file-manager/plugin/FileService.so")
	serviceManager.WaitForShutdown()
}
