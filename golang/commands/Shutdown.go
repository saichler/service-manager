package commands

import (
	. "github.com/saichler/console/golang/console"
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/common"
	"net"
	"strings"
)

type Shutdown struct {
	sm IServiceManager
}

func NewShutdown(sm IServiceManager) *Shutdown {
	sd := &Shutdown{}
	sd.sm = sm
	return sd
}

func (c *Shutdown) Command() string {
	return "shutdown"
}
func (c *Shutdown) Description() string {
	return "Shutdown the Service Manager"
}
func (c *Shutdown) Usage() string {
	return "shutdown"
}
func (c *Shutdown) ConsoleId() *ConsoleId {
	return c.sm.ConsoleId()
}
func (c *Shutdown) HandleCommand(command Command, args []string, conn net.Conn) (string, *ConsoleId) {
	Write("Are you sure you want to shutdown "+c.sm.ConsoleId().String()+" (yes/no)?", conn)
	reply, _ := Read(conn)
	reply = strings.ToLower(reply)
	for reply != "no" && reply != "yes" {
		Write("yes/no please?", conn)
		reply, _ := Read(conn)
		reply = strings.ToLower(reply)
	}
	if reply == "yes" {
		defer c.sm.Shutdown()
		return "Shutting Down...", nil
	}
	return "Canceled Shutdown.", nil
}
