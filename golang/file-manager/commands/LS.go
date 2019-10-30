package commands

import (
	"bytes"
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
)

type LS struct {
	service *FileManager
	mh      IMessageHandler
}

func NewLS(sm IService, mh IMessageHandler) *LS {
	sd := &LS{}
	sd.service = sm.(*FileManager)
	sd.mh = mh
	return sd
}

func (cmd *LS) Command() string {
	return "ls"
}

func (cmd *LS) Description() string {
	return "List the pre-set peer files at a pre-set location"
}

func (cmd *LS) Usage() string {
	return "ls"
}

func (cmd *LS) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *LS) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	response := cmd.mh.Request(cmd.service.PeerDir(), cmd.service.PeerServiceID())
	fd := response.(*model.FileDescriptor)
	buff := bytes.Buffer{}
	buff.WriteString("------------------------------------------------\n")
	buff.WriteString("Peer: " + cmd.service.PeerServiceID().String())
	buff.WriteString("\n")
	buff.WriteString("Directory: ")
	buff.WriteString(cmd.service.PeerDir())
	buff.WriteString("\n")

	for _, file := range fd.Files() {
		buff.WriteString("  - ")
		buff.WriteString(file.Name())
		buff.WriteString("\n")
	}
	return buff.String(), nil
}