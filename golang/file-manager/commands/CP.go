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
	"strings"
)

type CP struct {
	service *FileManagerService
	ls      IMessageHandler
	cp      IMessageHandler
	conn    net.Conn
}

func NewCpCMD(sm IService, rls, cp IMessageHandler) *CP {
	sd := &CP{}
	sd.service = sm.(*FileManagerService)
	sd.ls = rls
	sd.cp = cp
	if rls == nil {
		panic("Message handler is nil")
	}
	return sd
}

func (cmd *CP) Command() string {
	return "cp"
}

func (cmd *CP) Description() string {
	return "Copy a file from the remote location"
}

func (cmd *CP) Usage() string {
	return "cp <remote path> <local path>"
}

func (cmd *CP) ConsoleId() *ConsoleId {
	return cmd.service.ConsoleId()
}

func (cmd *CP) HandleCommand(command Command, args []string, conn net.Conn, id *ConsoleId) (string, *ConsoleId) {
	if len(args) < 2 {
		return cmd.Usage(), nil
	}
	rfilename := cmd.service.PeerDir() + "/" + args[0]
	lfilename := cmd.service.LocalDir() + "/" + args[1]

	req := model.NewFileRequest(rfilename, 1, true)
	response := cmd.ls.Request(req, cmd.service.PeerServiceID())
	fd := response.(*model.FileDescriptor)
	if fd.Name() == "" {
		return "File " + rfilename + " does not exit.", nil
	} else if fd.Size() == 0 {
		return "File in empty", nil
	}

	if _, err := os.Stat(lfilename); !os.IsNotExist(err) {
		hash, _ := security.FileHash(lfilename)
		if hash == fd.Hash() {
			return "File " + rfilename + " already exist in local dir", nil
		}
		console.Write("File already exist in local, overwrite (yes/no)?", conn)
		resp, _ := console.Read(conn)
		if resp != "yes" {
			return "Aborting", nil
		}
	}

	parts := fd.Parts()
	msg := lfilename + " (" + strconv.Itoa(int(fd.Size())) + "):"
	conn.Write([]byte(msg))
	if parts == 1 {
		data := cmd.cp.Request(fd, cmd.service.PeerServiceID()).(*model.FileData)
		ioutil.WriteFile(lfilename, data.Data(), 777)
	} else {
		cmd.conn = conn
		tasks := utils.NewJob(1, cmd)
		for i := 0; i < parts; i++ {
			fileData := model.NewFileData(rfilename, i, fd.Size())
			fpt := NewFetchPartTask(fileData, lfilename, cmd.cp, cmd.service)
			tasks.AddTask(fpt)
		}
		tasks.Run()
		assemble(lfilename)

	}
	hash, _ := security.FileHash(lfilename)
	valid := hash == fd.Hash()
	if valid {
		return "Done!", nil
	}
	return "Corrupted", nil
}

func (cmd *CP) Finished(task utils.JobTask) {
	if cmd.conn != nil {
		cmd.conn.Write([]byte("."))
	}
}

type FetchPartTask struct {
	fileData *model.FileData
	cp       IMessageHandler
	service  *FileManagerService
	target   string
}

func NewFetchPartTask(fileData *model.FileData, target string, cp IMessageHandler, service *FileManagerService) *FetchPartTask {
	fpt := &FetchPartTask{}
	fpt.fileData = fileData
	fpt.cp = cp
	fpt.service = service
	fpt.target = target
	return fpt
}

func (task *FetchPartTask) Run() {
	data := task.cp.Request(task.fileData, task.service.PeerServiceID()).(*model.FileData)
	file, _ := os.Create(task.target + ".part-" + getPartString(task.fileData.Part()))
	file.Write(data.Data())
	file.Close()
}

func getPartString(part int) string {
	str := strconv.Itoa(part)
	buff := bytes.Buffer{}
	for i := len(str); i < 5; i++ {
		buff.WriteString("0")
	}
	buff.WriteString(str)
	return buff.String()
}

func assemble(filename string) {
	os.Remove(filename)
	index := strings.LastIndex(filename, "/")
	dir := filename[0:index]
	filenames := make([]string, 0)
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		buff := bytes.Buffer{}
		buff.WriteString(dir)
		buff.WriteString("/")
		buff.WriteString(f.Name())
		fp := buff.String()
		if len(fp) < len(filename) {
			continue
		}
		if fp[0:len(filename)] == filename {
			filenames = append(filenames, fp)
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
	file.Close()
	file, _ = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	for _, fn := range filenames {
		data, _ := ioutil.ReadFile(fn)
		file.Write(data)
		os.Remove(fn)
	}
	file.Close()
}
