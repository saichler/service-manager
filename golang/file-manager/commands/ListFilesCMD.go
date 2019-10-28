package commands

import (
	"bytes"
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
)

type ListFiledCMD struct {
	service *FileManager
	mh      IMessageHandler
}

func NewListFiels(sm IService, mh IMessageHandler) *ListFiledCMD {
	sd := &ListFiledCMD{}
	sd.service = sm.(*FileManager)
	sd.mh = mh
	if mh == nil {
		panic("Message handler is nil")
	}
	return sd
}

func (c *ListFiledCMD) Command() string {
	return "list-files"
}

func (c *ListFiledCMD) Description() string {
	return "List a peer files at a location"
}

func (c *ListFiledCMD) Usage() string {
	return "list-files <path>"
}

func (c *ListFiledCMD) ConsoleId() *ConsoleId {
	return c.service.ConsoleId()
}

func (c *ListFiledCMD) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	dest := &protocol.ServiceID{}
	sid := args[0] + " " + args[1]
	e := dest.Parse(sid)
	if e != nil {
		return "Invalid service id format:" + args[0], nil
	}
	response := c.mh.Request("/tmp", dest)
	fd := response.(*model.FileDescriptor)
	buff := bytes.Buffer{}
	for _, file := range fd.Files() {
		buff.WriteString(file.Name())
		buff.WriteString("\n")
	}
	return buff.String(), nil
}
