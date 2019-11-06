package commands

import (
	"bytes"
	"github.com/saichler/console/golang/console"
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"net"
)

type Diff struct {
	service *FileManagerService
	ls      IMessageHandler
}

func NewDiff(sm IService, ls IMessageHandler) *Diff {
	sd := &Diff{}
	sd.service = sm.(*FileManagerService)
	sd.ls = ls
	return sd
}

func (cmd *Diff) Command() string {
	return "diff"
}

func (cmd *Diff) Description() string {
	return "Compare remote and local dirs"
}

func (cmd *Diff) Usage() string {
	return "diff"
}

func (cmd *Diff) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *Diff) HandleCommand(args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	dirRequest := model.NewFileRequest(cmd.service.PeerDir(), 100, false)
	console.Write("Scanning Remote Directory, this may take a min...", conn)
	response := cmd.ls.Request(dirRequest, cmd.service.PeerServiceID())
	console.Writeln("Done!", conn)
	aside := response.(*model.FileDescriptor)
	zside := model.NewFileDescriptor(cmd.service.LocalDir()+"/"+aside.Name(), 100, false)
	aSideMissing, zSideMissing := diff(aside, zside)
	buff := bytes.Buffer{}
	buff.WriteString("Remote missing files:\n")
	for name, _ := range aSideMissing {
		buff.WriteString(name)
		buff.WriteString("\n")
	}
	buff.WriteString("Local missing files:\n")
	for name, _ := range zSideMissing {
		buff.WriteString(name)
		buff.WriteString("\n")
	}
	return buff.String(), nil
}

func diff(aside, zside *model.FileDescriptor) (map[string]string, map[string]string) {
	zsideTarget := model.NewFileDescriptor(aside.SourceParent().SourcePath(), 0, false)
	asideTarget := model.NewFileDescriptor(zside.SourceParent().SourcePath(), 0, false)
	aside.SetTargetParent(zsideTarget)
	zside.SetTargetParent(asideTarget)
	aSideMissing := make(map[string]string)
	zSideMissing := make(map[string]string)
	deepDiff(aside, zside, aside.SourceRoot(), zside.SourceRoot(), aSideMissing, zSideMissing)
	return aSideMissing, zSideMissing
}

func deepDiff(aside, zside, aSideRoot, zSideRoot *model.FileDescriptor, aSideMissing, zSideMissing map[string]string) {
	if aside == nil && zside != nil {
		aSideMissing[zside.SourcePath()] = zside.SourcePath()
		return
	}
	if aside != nil && zside == nil {
		zSideMissing[aside.SourcePath()] = aside.SourcePath()
		return
	}

	if aside.IsDir() {
		for _, aSideChild := range aside.Files() {
			path := aSideChild.TargetPath()
			zSideChild := zSideRoot.Get(path)
			deepDiff(aSideChild, zSideChild, aSideRoot, zSideRoot, aSideMissing, zSideMissing)
		}
	}

	if zside.IsDir() {
		for _, zSideChild := range zside.Files() {
			path := zSideChild.TargetPath()
			aSideChild := aSideRoot.Get(path)
			deepDiff(aSideChild, zSideChild, aSideRoot, zSideRoot, aSideMissing, zSideMissing)
		}
	}
}
