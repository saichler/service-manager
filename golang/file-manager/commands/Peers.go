package commands

import (
	"bytes"
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
)

type Peers struct {
	service *FileManagerService
}

func NewPeers(sm IService) *Peers {
	sd := &Peers{}
	sd.service = sm.(*FileManagerService)
	return sd
}

func (c *Peers) Command() string {
	return "peers"
}

func (c *Peers) Description() string {
	return "List all the file manager peers."
}

func (c *Peers) Usage() string {
	return "peers"
}

func (c *Peers) ConsoleId() *ConsoleId {
	return c.service.ConsoleId()
}

func (c *Peers) HandleCommand(args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	peers := c.service.ServiceManager().ServiceNetwork().GetPeers(c.service.ServiceID())
	buff := bytes.Buffer{}
	for _, peer := range peers {
		buff.WriteString(peer.String())
		buff.WriteString("\n")
	}
	if len(peers) == 0 {
		buff.WriteString("No Peers")
	}
	return buff.String(), nil
}
