package main

import (
	"github.com/saichler/service-manager/golang/file-manager/commands"
	message_handlers "github.com/saichler/service-manager/golang/file-manager/message-handlers"
	"github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

var Service IService = &service.FileManager{}
var Commands IServiceCommands = &commands.FileManagerCommands{}
var Handlers IServiceMessageHandlers = &message_handlers.FileManagerHandlers{}
