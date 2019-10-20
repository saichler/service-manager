package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/common"
	"net"
	"strconv"
	"strings"
)

type CD struct {
	sm IServiceManager
	sv IService
}

func NewCD(context interface{}) *CD {
	sd := &CD{}
	sm, ok := context.(IServiceManager)
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

func (c *CD) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	if len(args) == 0 {
		return "Service name is required", nil
	}
	if c.sm != nil {
		return c.handleServiceManager(command, args, conn, id)
	} else {
		return c.handleService(command, args, conn, id)
	}
}

func (c *CD) handleServiceManager(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	if len(args) == 0 {
		return "Service name is required", nil
	}
	serviceID := args[0]
	if len(args) > 1 {
		for i := 1; i < len(args); i++ {
			serviceID += " " + args[i]
		}
	}
	serviceID = strings.ToLower(serviceID)
	services := c.sm.Services()
	for _, service := range services {
		id := strings.ToLower(service.Topic() + "-" + strconv.Itoa(int(service.ID())))
		if id == serviceID {
			return "", service.ConsoleId()
		}
	}
	return "Unknown service " + serviceID, nil
}

func (c *CD) handleService(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	if args[0] == ".." && id.Parent() != nil {
		return "", id.Parent()
	}
	return "", nil
}
