package commands

import (
	"github.com/saichler/console/golang/console"
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
	"strconv"
)

type Sync struct {
	service *FileManagerService
	cp      IMessageHandler
	ls      IMessageHandler
}

func NewSync(sm IService, cp IMessageHandler, ls IMessageHandler) *Sync {
	sd := &Sync{}
	sd.service = sm.(*FileManagerService)
	sd.cp = cp
	sd.ls = ls
	return sd
}

func (cmd *Sync) Command() string {
	return "sync"
}

func (cmd *Sync) Description() string {
	return "Sync remote direcory to local."
}

func (cmd *Sync) Usage() string {
	return "sync"
}

func (cmd *Sync) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *Sync) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	fr := model.NewFileRequest(cmd.service.PeerDir(), 100)
	console.Write("Calculating Hashes, this may take a min...", conn)
	response := cmd.ls.Request(fr, cmd.service.PeerServiceID())
	console.Writeln("Done!", conn)
	fd := response.(*model.FileDescriptor)
	files, dirs := countFilesAndDirectories(fd)
	msg := "Going to sync " + strconv.Itoa(files) + " files in " + strconv.Itoa(dirs) + " directories."
	return msg, nil
}

func countFilesAndDirectories(fileDescriptor *model.FileDescriptor) (int, int) {
	if fileDescriptor.Files() == nil || len(fileDescriptor.Files()) == 0 {
		return 1, 0
	}
	dirs := 1
	files := 0
	for _, child := range fileDescriptor.Files() {
		f, d := countFilesAndDirectories(child)
		dirs += d
		files += f
	}
	return files, dirs
}
