package commands

import (
	"github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/common"
	"net"
)

type Uplink struct {
	ms IService
}

func NewUplink(ms IService) *Uplink {
	sd := &Uplink{}
	sd.ms = ms
	return sd
}

func (c *Uplink) Command() string {
	return "uplink"
}
func (c *Uplink) Description() string {
	return "creates an uplink to another node"
}
func (c *Uplink) Usage() string {
	return "uplink"
}
func (c *Uplink) ConsoleId() *commands.ConsoleId {
	return c.ms.ConsoleId()
}
func (c *Uplink) HandleCommand(command commands.Command, args []string, conn net.Conn, id *commands.ConsoleId) (string, *commands.ConsoleId) {
	if len(args) == 0 {
		return "Uplink require a dest ip", nil
	}
	c.ms.ServiceManager().Uplink(args[0])
	return "Sent a message to " + args[0], nil
}
