package main

import (
	"github.com/saichler/service-manager/golang/file-manager/commands"
	message_handlers "github.com/saichler/service-manager/golang/file-manager/message-handlers"
	"github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

var Service, Commands, Handlers = NewService()

func NewService() (IService, IServiceCommands, IServiceMessageHandlers) {
	s := &service.FileManager{}
	h := &message_handlers.FileManagerHandlers{}
	h.Init(s)
	c := &commands.FileManagerCommands{}
	c.Init(s, h)
	return s, c, h
}
