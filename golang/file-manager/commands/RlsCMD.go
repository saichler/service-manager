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

type RlsCMD struct {
	service *FileManager
	mh      IMessageHandler
}

func NewRlsCMD(sm IService, mh IMessageHandler) *RlsCMD {
	sd := &RlsCMD{}
	sd.service = sm.(*FileManager)
	sd.mh = mh
	if mh == nil {
		panic("Message handler is nil")
	}
	return sd
}

func (c *RlsCMD) Command() string {
	return "rls"
}

func (c *RlsCMD) Description() string {
	return "List a peer files at a location"
}

func (c *RlsCMD) Usage() string {
	return "rls <serviceid> <path>"
}

func (c *RlsCMD) ConsoleId() *ConsoleId {
	return c.service.ConsoleId()
}

func (c *RlsCMD) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	dest := &protocol.ServiceID{}
	if len(args) < 2 {
		return c.Usage(), nil
	}
	e := dest.Parse(args[0])
	if e != nil {
		return "Invalid service id format:" + args[0], nil
	}
	response := c.mh.Request(args[1], dest)
	fd := response.(*model.FileDescriptor)
	buff := bytes.Buffer{}
	for _, file := range fd.Files() {
		buff.WriteString(file.Name())
		buff.WriteString("\n")
	}
	return buff.String(), nil
}
