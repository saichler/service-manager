package commands

import (
	"github.com/saichler/console/golang/console/commands"
	commands2 "github.com/saichler/service-manager/golang/management/commands"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type FileManagerCommands struct {
	commands map[string]commands.Command
}

func (c *FileManagerCommands) Init(service IService, mh IServiceMessageHandlers) {
	c.commands = make(map[string]commands.Command)
	c.addCommand(commands2.NewCD(service))
	c.addCommand(NewListPeers(service))
	c.addCommand(NewRlsCMD(service, mh.Handler("RLS")))
	c.addCommand(NewCpCMD(service, mh.Handler("RLS")))
}

func (c *FileManagerCommands) addCommand(cmd commands.Command) {
	c.commands[cmd.Command()] = cmd
}

func (c *FileManagerCommands) Commands() []commands.Command {
	result := make([]commands.Command, 0)
	for _, cmd := range c.commands {
		result = append(result, cmd)
	}
	return result
}
