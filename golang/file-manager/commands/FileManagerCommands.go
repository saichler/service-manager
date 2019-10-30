package commands

import (
	"github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type FileManagerCommands struct {
	commands map[string]commands.Command
}

func (c *FileManagerCommands) Init(service IService, mh IServiceMessageHandlers) {
	c.commands = make(map[string]commands.Command)
	c.addCommand(NewListPeers(service))
	c.addCommand(NewLS(service, mh.Handler("ls")))
	c.addCommand(NewCD(service, mh.Handler("ls")))
	c.addCommand(NewCpCMD(service, mh.Handler("ls")))
	//c.addAlias(commands2.NewCD(service), "lcd")
}

func (c *FileManagerCommands) addCommand(cmd commands.Command) {
	c.commands[cmd.Command()] = cmd
}

func (c *FileManagerCommands) addAlias(cmd commands.Command, alias string) {
	c.commands[alias] = cmd
}

func (c *FileManagerCommands) Commands() []commands.Command {
	result := make([]commands.Command, 0)
	for _, cmd := range c.commands {
		result = append(result, cmd)
	}
	return result
}
