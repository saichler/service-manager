package commands

import (
	"github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type ManagementCommands struct {
	commands map[string]commands.Command
}

func (c *ManagementCommands) Init(service IService, mh IServiceMessageHandlers) {
	c.commands = make(map[string]commands.Command)
	c.addCommand(NewShutdown(service))
	c.addCommand(NewUplink(service))
	c.addCommand(NewLoadService(service))
	c.addCommand(NewCD(service))

}

func (c *ManagementCommands) addCommand(cmd commands.Command) {
	c.commands[cmd.Command()] = cmd
}

func (c *ManagementCommands) Commands() map[string]commands.Command {
	return c.commands
}
