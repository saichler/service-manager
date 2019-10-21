package commands

import (
	"bytes"
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/common"
	"net"
)

type LS struct {
	sm IServiceManager
}

func NewLS(sm IServiceManager) *LS {
	sd := &LS{}
	sd.sm = sm
	return sd
}

func (c *LS) Command() string {
	return "ls"
}
func (c *LS) Description() string {
	return "list the services in the system"
}
func (c *LS) Usage() string {
	return "ls"
}
func (c *LS) ConsoleId() *ConsoleId {
	return c.sm.ConsoleId()
}
func (c *LS) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	buff := bytes.Buffer{}
	services := c.sm.Services()
	for _, service := range services {
		buff.WriteString(service.ConsoleId().Key())
		buff.WriteString("\n")
	}
	return buff.String(), nil
}