package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
	"strconv"
)

type CpCMD struct {
	service *FileManager
	rls     IMessageHandler
}

func NewCpCMD(sm IService, rls IMessageHandler) *CpCMD {
	sd := &CpCMD{}
	sd.service = sm.(*FileManager)
	sd.rls = rls
	if rls == nil {
		panic("Message handler is nil")
	}
	return sd
}

func (c *CpCMD) Command() string {
	return "cp"
}

func (c *CpCMD) Description() string {
	return "Copy a file from the remote location"
}

func (c *CpCMD) Usage() string {
	return "cp <path remote> <path local>"
}

func (c *CpCMD) ConsoleId() *ConsoleId {
	return c.service.ConsoleId()
}

func (c *CpCMD) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	dest := &protocol.ServiceID{}
	if len(args) < 2 {
		return c.Usage(), nil
	}
	e := dest.Parse(args[0])
	if e != nil {
		return "Invalid service id format:" + args[0], nil
	}
	response := c.rls.Request(args[1], dest)
	fd := response.(*model.FileDescriptor)

	return "remote size:" + strconv.Itoa(int(fd.Size())), nil
}
