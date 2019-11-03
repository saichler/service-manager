package commands

import (
	"fmt"
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
)

type Set struct {
	service *FileManagerService
}

func NewSet(sm IService) *Set {
	sd := &Set{}
	sd.service = sm.(*FileManagerService)
	return sd
}

func (c *Set) Command() string {
	return "set"
}

func (c *Set) Description() string {
	return "set the active peer"
}

func (c *Set) Usage() string {
	return "set <peer id>"
}

func (c *Set) ConsoleId() *ConsoleId {
	return c.service.ConsoleId()
}

func (c *Set) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	if len(args) == 0 {
		return c.Usage(), nil
	}
	peer := &protocol.ServiceID{}
	fmt.Println(args[0])
	e := peer.Parse(args[0])
	if e != nil {
		return "Invalid peer id: " + args[0] + ", did you forget the '?", nil
	}
	c.service.SetPeerServiceID(peer)

	peers := c.service.ServiceManager().ServiceNetwork().GetPeers(c.service.ServiceID())
	found := false
	for _, p := range peers {
		if p.String() == peer.String() {
			found = true
			break
		}
	}
	msg := "Peer set to:" + peer.String()
	if !found {
		msg += " but not found!"
	}
	return msg, nil
}
