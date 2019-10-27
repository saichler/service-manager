package commands

import (
	"github.com/saichler/console/golang/console/commands"
	commands2 "github.com/saichler/service-manager/golang/management/commands"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type FileManagerCommands struct {
}

func (mcmd *FileManagerCommands) Commands(service IService) []commands.Command {
	result := make([]commands.Command, 0)
	result = append(result, commands2.NewCD(service))
	result = append(result, NewListPeers(service))
	return result
}
