package main

import (
	commands2 "github.com/saichler/service-manager/golang/file-manager/commands"
	message_handlers2 "github.com/saichler/service-manager/golang/file-manager/message-handlers"
	service2 "github.com/saichler/service-manager/golang/file-manager/service"
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

	serviceManager.Console().RegisterCommand(commands.NewLS(serviceManager), "ls")
	serviceManager.Console().RegisterCommand(commands.NewCD(serviceManager), "cd")
	serviceManager.AddService(NewService())
	serviceManager.AddService(NewFileService())
	/*
		files, e := ioutil.ReadDir("./plugins")
		if e == nil {
			for _, f := range files {
				if strings.Contains(f.Name(), ".so") {
					serviceManager.LoadService("./plugins/" + f.Name())
				}
			}
		} else {
			os.Create("./plugins")
		}*/
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

func NewFileService() (service_manager.IService, service_manager.IServiceCommands, service_manager.IServiceMessageHandlers) {
	s := &service2.FileManager{}
	h := &message_handlers2.FileManagerHandlers{}
	h.Init(s)
	c := &commands2.FileManagerCommands{}
	c.Init(s, h)
	return s, c, h
}
