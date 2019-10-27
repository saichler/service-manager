package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
)

type ListPeers struct {
	service *FileManager
}

func NewListPeers(sm IService) *ListPeers {
	sd := &ListPeers{}
	sd.service = sm.(*FileManager)
	return sd
}

func (c *ListPeers) Command() string {
	return "list-peers"
}

func (c *ListPeers) Description() string {
	return "List all the file manager peers."
}

func (c *ListPeers) Usage() string {
	return "list-peers"
}

func (c *ListPeers) ConsoleId() *ConsoleId {
	return c.service.ConsoleId()
}

func (c *ListPeers) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	c.service.ServiceManager().ServiceNetwork().
	return "", nil
}
