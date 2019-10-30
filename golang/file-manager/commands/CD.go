package commands

import (
	. "github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
)

type CD struct {
	service *FileManager
	mh      IMessageHandler
}

func NewCD(sm IService, mh IMessageHandler) *CD {
	sd := &CD{}
	sd.service = sm.(*FileManager)
	sd.mh = mh
	return sd
}

func (cmd *CD) Command() string {
	return "cd"
}

func (cmd *CD) Description() string {
	return "Change directory"
}

func (cmd *CD) Usage() string {
	return "cd <dir>"
}

func (cmd *CD) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *CD) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {

	if len(args) == 0 {
		return cmd.Usage(), nil
	}
	cmd.service.SetPeerDir(args[0])

	return "", nil
}
