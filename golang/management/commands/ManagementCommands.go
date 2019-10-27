package commands

import (
	"github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type ManagementCommands struct {
}

func (mcmd *ManagementCommands) Commands(service IService) []commands.Command {
	result := make([]commands.Command, 0)
	result = append(result, NewShutdown(service))
	result = append(result, NewUplink(service))
	result = append(result, NewLoadService(service))
	result = append(result, NewCD(service))
	return result
}
