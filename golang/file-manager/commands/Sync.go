package commands

import (
	"bytes"
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
	"sort"
	"strconv"
	"time"
)

type Sync struct {
	service     *FileManagerService
	cp          IMessageHandler
	ls          IMessageHandler
	conn        net.Conn
	running     bool
	currentSize int64
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

func (cmd *Sync) HandleCommand(args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	cmd.running = true
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
	if !cmd.running {
		return "", nil
	}
	report := sr.Report(true)
	console.Writeln(report, cmd.conn)
	console.Writeln("Start downloading size diff files...", conn)
	sideDiff := sr.SizeDiff()
	sr = model.NewSyncReport()
	for _, d := range sideDiff {
		cmd.copyFile(d, sr)
		if !cmd.running {
			return "", nil
		}
	}
	report = sr.Report(true)
	console.Writeln(report, cmd.conn)
	return "Done", nil
}

func (cmd *Sync) copyFiles(descriptor *model.FileDescriptor, sr *model.SyncReport) {
	if descriptor.IsDir() {
		for _, file := range descriptor.Files() {
			cmd.copyFiles(file, sr)
			if !cmd.running {
				return
			}
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

	if !cmd.running {
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

func (cmd *Sync) copySmallFile(descriptor *model.FileDescriptor) {
	fileData := model.NewFileData(descriptor.SourcePath(), 0, descriptor.Size())
	data := cmd.cp.Request(fileData, cmd.service.PeerServiceID()).(*model.FileData)
	ioutil.WriteFile(descriptor.TargetPath(), data.Data(), 777)
}

func (cmd *Sync) copyFile(descriptor *model.FileDescriptor, sr *model.SyncReport) {
	if !cmd.running {
		return
	}
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
		go cmd.copySmallFile(descriptor)
	} else {
		tasks := utils.NewJob(5, newCopyFileJobListener(descriptor.TargetPath(), parts, cmd))
		for i := 0; i < parts; i++ {
			fileData := model.NewFileData(descriptor.SourcePath(), i, descriptor.Size())
			fpt := NewFetchPartTask(fileData, descriptor.TargetPath(), cmd.cp, cmd.service)
			tasks.AddTask(fpt)
		}
		tasks.Run()
		assembleFile(descriptor)
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

func assembleFile(descriptor *model.FileDescriptor) {
	parent := descriptor.TargetParent()
	if parent == nil {
		parent = descriptor.SourceParent()
	}

	dir := parent.TargetPath()
	filename := descriptor.TargetPath()

	filenames := make([]string, 0)
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		buff := bytes.Buffer{}
		buff.WriteString(dir)
		buff.WriteString("/")
		buff.WriteString(f.Name())
		part := buff.String()
		if isPartOfFile(part, filename) {
			filenames = append(filenames, part)
		}
	}

	sort.Slice(filenames, func(i, j int) bool {
		attValueA := filenames[i]
		attValueB := filenames[j]
		if attValueA < attValueB {
			return true
		}
		return false
	})

	file, _ := os.Create(filename)
	for _, fn := range filenames {
		data, _ := ioutil.ReadFile(fn)
		file.Write(data)
		os.Remove(fn)
	}
	file.Close()
}

func isPartOfFile(f, filename string) bool {
	if len(f) > len(filename) && f[0:len(filename)] == filename {
		return true
	}
	return false
}

func (cmd *Sync) Join(conn net.Conn) {
	cmd.conn = conn
}

func (cmd *Sync) Stop(conn net.Conn) {
	cmd.running = false
	console.Writeln("Stop signal was sent", conn)
}

type CopyFileJobListener struct {
	filename    string
	cmd         *Sync
	finishCount int
	parts       int
	lastReport  int64
}

func newCopyFileJobListener(filename string, parts int, cmd *Sync) *CopyFileJobListener {
	cpjl := &CopyFileJobListener{}
	cpjl.filename = filename
	cpjl.parts = parts
	cpjl.cmd = cmd
	cpjl.lastReport = time.Now().Unix()
	return cpjl
}

func (jl *CopyFileJobListener) Finished(task utils.JobTask) {
	jl.finishCount++
	if jl.cmd.conn != nil {
		if time.Now().Unix()-jl.lastReport > 30 {
			p := float64(jl.finishCount) / float64(jl.parts) * 100
			console.Writeln("\n"+jl.filename+"("+strconv.Itoa(int(p))+"%).", jl.cmd.conn)
			jl.lastReport = time.Now().Unix()
		} else {
			console.Write(".", jl.cmd.conn)
		}
	}
}
