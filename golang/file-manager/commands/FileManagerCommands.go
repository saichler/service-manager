package commands

import (
	"github.com/saichler/console/golang/console/commands"
	message_handlers "github.com/saichler/service-manager/golang/file-manager/message-handlers"
	commands2 "github.com/saichler/service-manager/golang/management/commands"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type FileManagerCommands struct {
}

func (mcmd *FileManagerCommands) Commands(service IService) []commands.Command {
	result := make([]commands.Command, 0)
	result = append(result, commands2.NewCD(service))
	result = append(result, NewListPeers(service))
	result = append(result, NewListFiels(service, &message_handlers.ListFiles{}))
	return result
}
