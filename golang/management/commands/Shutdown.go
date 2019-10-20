package commands

import (
	. "github.com/saichler/console/golang/console"
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/common"
	"net"
	"strings"
)

type Shutdown struct {
	service IService
}

func NewShutdown(sm IService) *Shutdown {
	sd := &Shutdown{}
	sd.service = sm
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
	return c.service.ConsoleId()
}
func (c *Shutdown) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	Write("Are you sure you want to shutdown "+c.service.ServiceManager().ConsoleId().String()+" (yes/no)?", conn)
	reply, _ := Read(conn)
	reply = strings.ToLower(reply)
	for reply != "no" && reply != "yes" {
		Write("yes/no please?", conn)
		reply, _ := Read(conn)
		reply = strings.ToLower(reply)
	}
	if reply == "yes" {
		c.service.ServiceManager().Publish("Shutdown", c.service, []byte("Shutdown"))
		defer c.service.ServiceManager().Shutdown()
		return "Shutting Down...", nil
	}
	return "Canceled Shutdown.", nil
}
