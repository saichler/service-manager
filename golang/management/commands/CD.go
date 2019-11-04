package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
	"strconv"
	"strings"
)

type CD struct {
	sm *ServiceManager
	sv IService
}

func NewCD(context interface{}) *CD {
	sd := &CD{}
	sm, ok := context.(*ServiceManager)
	if ok {
		sd.sm = sm
	} else {
		sd.sv = context.(IService)
	}
	return sd
}
func (c *CD) Command() string {
	return "cd"
}
func (c *CD) Description() string {
	return "Go in & out of a service entry"
}
func (c *CD) Usage() string {
	return "cd <service instance>/.."
}
func (c *CD) ConsoleId() *ConsoleId {
	if c.sm != nil {
		return c.sm.ConsoleId()
	}
	return c.sv.ConsoleId()
}

func (c *CD) HandleCommand(args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	if len(args) == 0 {
		if id.Parent() != nil {
			return "", id.Parent()
		}
		return "Service name is required", nil
	}
	if c.sm != nil {
		return c.handleServiceManager(c, args, conn, id)
	} else {
		return c.handleService(c, args, conn, id)
	}
}

func (c *CD) handleServiceManager(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	serviceID := args[0]
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			serviceID += " " + args[i]
		}
	}
	serviceID = strings.ToLower(serviceID)
	services := c.sm.Services()
	var subset IService
	for _, service := range services {
		id := strings.ToLower(service.Topic() + "-" + strconv.Itoa(int(service.ID())))
		if id == serviceID {
			return "", service.ConsoleId()
		}
		if len(serviceID) < len(id) && id[0:len(serviceID)] == serviceID {
			subset = service
		}
	}
	if subset != nil {
		return "", subset.ConsoleId()
	}
	return "Unknown service " + serviceID, nil
}

func (c *CD) handleService(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	if args[0] == ".." && id.Parent() != nil {
		return "", id.Parent()
	}
	return "", nil
}
