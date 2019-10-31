package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

type CpCMD struct {
	service *FileManager
	ls      IMessageHandler
	cp      IMessageHandler
}

func NewCpCMD(sm IService, rls, cp IMessageHandler) *CpCMD {
	sd := &CpCMD{}
	sd.service = sm.(*FileManager)
	sd.ls = rls
	sd.cp = cp
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
	return "cp <remote path> <local path>"
}

func (c *CpCMD) ConsoleId() *ConsoleId {
	return c.service.ConsoleId()
}

func (c *CpCMD) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	if len(args) < 2 {
		return c.Usage(), nil
	}
	rfilename := c.service.PeerDir() + "/" + args[0]
	lfilename := c.service.LocalDir() + "/" + args[1]
	response := c.ls.Request(rfilename, c.service.PeerServiceID())
	fd := response.(*model.FileDescriptor)
	if fd.Name() == "" {
		return "File " + rfilename + " does not exit.", nil
	}

	parts := fd.Part()

	if parts == 1 {
		fd.SetPart(0)
		data := c.cp.Request(fd, c.service.PeerServiceID()).(*model.FileData)
		ioutil.WriteFile(lfilename, data.Data(), 777)
		return "Written file:" + lfilename + " size:" + strconv.Itoa(int(fd.Size())), nil
	} else {
		file, _ := os.Create(lfilename)
		file.Close()
		file, _ = os.OpenFile(lfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		for i := 0; i < parts; i++ {
			fd.SetPart(i)
			data := c.cp.Request(fd, c.service.PeerServiceID()).(*model.FileData)
			file.Write(data.Data())
		}
		file.Close()
	}
	return "Multipart File with:" + strconv.Itoa(parts) + " and the size of:" + strconv.Itoa(int(fd.Size())), nil
}
