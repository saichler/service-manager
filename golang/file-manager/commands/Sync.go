package commands

import (
	"github.com/saichler/console/golang/console"
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/security"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	utils "github.com/saichler/utils/golang"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

type Sync struct {
	service *FileManagerService
	cp      IMessageHandler
	ls      IMessageHandler
	conn    net.Conn
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
	dirRequest := model.NewFileRequest(cmd.service.PeerDir(), 100, false)
	console.Write("Scanning Remote Directory, this may take a min...", conn)
	response := cmd.ls.Request(dirRequest, cmd.service.PeerServiceID())
	console.Writeln("Done!", conn)
	descriptor := response.(*model.FileDescriptor)
	files, dirs := countFilesAndDirectories(descriptor)
	msg := "Going to sync " + strconv.Itoa(files) + " files in " + strconv.Itoa(dirs) + " directories."
	console.Writeln(msg, conn)
	msg = "Creating " + strconv.Itoa(dirs) + " on local target..."
	console.Write(msg, conn)
	cmd.createDirectories(descriptor)
	console.Writeln("Done!", conn)
	console.Writeln("Start downloading new files...", conn)
	cmd.conn = conn
	sr := model.NewSyncReport()
	cmd.copyFiles(descriptor, sr)
	report := sr.Report(true)
	console.Writeln(report, cmd.conn)
	console.Writeln("Start downloading size diff files...", conn)
	sideDiff := sr.SizeDiff()
	sr = model.NewSyncReport()
	for _, d := range sideDiff {
		cmd.copyFile(d, sr)
	}
	report = sr.Report(true)
	console.Writeln(report, cmd.conn)
	return "Done", nil
}

func (cmd *Sync) copyFiles(descriptor *model.FileDescriptor, sr *model.SyncReport) {
	if descriptor.IsDir() {
		for _, file := range descriptor.Files() {
			cmd.copyFiles(file, sr)
		}
		return
	}

	if descriptor.Name() == "" {
		sr.AddErrored(descriptor)
		return
	} else if descriptor.Size() == 0 {
		sr.AddErrored(descriptor)
		return
	}

	exists, err := os.Stat(descriptor.TargetPath())

	if !os.IsNotExist(err) {
		if descriptor.Size() == exists.Size() {
			sr.AddExist(descriptor)
			return
		} else {
			sr.AddSizeDiff(descriptor)
			return
		}
		/*
			request := model.NewFileRequest(descriptor.SourcePath(), 1, true)
			response := cmd.ls.Request(request, cmd.service.PeerServiceID())
			descriptor.SetHash(response.(*model.FileDescriptor).Hash())
			hash, _ := security.FileHash(descriptor.TargetPath())
			if hash == descriptor.Hash() {
				return
			}*/
	}

	cmd.copyFile(descriptor, sr)
	report := sr.Report(false)
	if report != "" {
		console.Writeln(report, cmd.conn)
	}
}

func (cmd *Sync) Finished(task utils.JobTask) {
	if cmd.conn != nil {
		cmd.conn.Write([]byte("."))
	}
}

func (cmd *Sync) copyFile(descriptor *model.FileDescriptor, sr *model.SyncReport) {
	msg := descriptor.TargetPath() + " (" + strconv.Itoa(int(descriptor.Size())) + "): "
	console.Write(msg, cmd.conn)

	if _, err := os.Stat(descriptor.TargetPath()); !os.IsNotExist(err) {
		hash, _ := security.FileHash(descriptor.TargetPath())
		if hash == descriptor.Hash() {
			sr.AddExist(descriptor)
			return
		}
	}

	parts := descriptor.Parts()
	if parts == 1 {
		fileData := model.NewFileData(descriptor.SourcePath(), 0, descriptor.Size())
		data := cmd.cp.Request(fileData, cmd.service.PeerServiceID()).(*model.FileData)
		ioutil.WriteFile(descriptor.TargetPath(), data.Data(), 777)
	} else {
		tasks := utils.NewJob(1, cmd)
		for i := 0; i < parts; i++ {
			fileData := model.NewFileData(descriptor.SourcePath(), i, descriptor.Size())
			fpt := NewFetchPartTask(fileData, descriptor.TargetPath(), cmd.cp, cmd.service)
			tasks.AddTask(fpt)
		}
		tasks.Run()
		assemble(descriptor.TargetPath())
	}

	sr.AddCopied(descriptor)

	if descriptor.Hash() != "" {
		hash, _ := security.FileHash(descriptor.TargetPath())
		valid := hash == descriptor.Hash()
		if valid {
			console.Writeln("Done!", cmd.conn)
		} else {
			console.Writeln("Corrupted!", cmd.conn)
		}
	} else {
		console.Writeln("Done!", cmd.conn)
	}
}

func (cmd *Sync) createDirectories(descriptor *model.FileDescriptor) {
	local := model.NewFileDescriptor(cmd.service.LocalDir(), 0, false)
	descriptor.SetTargetParent(local)
	createDirectories(descriptor)
}

func createDirectories(descriptor *model.FileDescriptor) {
	if descriptor.Files() == nil || len(descriptor.Files()) == 0 {
		return
	}
	_, e := os.Stat(descriptor.TargetPath())
	if os.IsNotExist(e) {
		os.MkdirAll(descriptor.TargetPath(), 0777)
	}
	for _, child := range descriptor.Files() {
		createDirectories(child)
	}
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
