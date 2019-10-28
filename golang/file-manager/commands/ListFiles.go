package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/messaging/golang/net/protocol"
	message_handlers "github.com/saichler/service-manager/golang/file-manager/message-handlers"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
)

type ListFiles struct {
	service *FileManager
	mh      *message_handlers.ListFiles
}

func NewListFiels(sm IService, files *message_handlers.ListFiles) *ListFiles {
	sd := &ListFiles{}
	sd.service = sm.(*FileManager)
	return sd
}

func (c *ListFiles) Command() string {
	return "list-files"
}

func (c *ListFiles) Description() string {
	return "List a peer files at a location"
}

func (c *ListFiles) Usage() string {
	return "list-files <path>"
}

func (c *ListFiles) ConsoleId() *ConsoleId {
	return c.service.ConsoleId()
}

func (c *ListFiles) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	dest := &protocol.ServiceID{}
	e := dest.Parse(args[0])
	if e != nil {
		return "Invalid service id format:" + args[0], nil
	}
	e = c.service.ServiceManager().Send(c.mh.Topic(), c.service, dest, []byte("/tmp"), false)
	if e != nil {
		return e.Error(), nil
	}
}
