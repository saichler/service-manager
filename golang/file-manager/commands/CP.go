package commands

import (
	"bytes"
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
	service *FileManager
	ls      IMessageHandler
	cp      IMessageHandler
	conn    net.Conn
}

func NewCpCMD(sm IService, rls, cp IMessageHandler) *CP {
	sd := &CP{}
	sd.service = sm.(*FileManager)
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
	response := cmd.ls.Request(rfilename, cmd.service.PeerServiceID())
	fd := response.(*model.FileDescriptor)
	if fd.Name() == "" {
		return "File " + rfilename + " does not exit.", nil
	}

	parts := fd.Part()
	msg := "(" + strconv.Itoa(int(fd.Size())) + ")"
	conn.Write([]byte(msg))
	if parts == 1 {
		fd.SetPart(0)
		data := cmd.cp.Request(fd, cmd.service.PeerServiceID()).(*model.FileData)
		ioutil.WriteFile(lfilename, data.Data(), 777)
	} else {
		cmd.conn = conn
		tasks := utils.NewJob(5, cmd)
		for i := 0; i < parts; i++ {
			descriptor := fd.Clone()
			descriptor.SetPart(i)
			fpt := NewFetchPartTask(descriptor, lfilename, cmd.cp, cmd.service)
			tasks.AddTask(fpt)
		}
		tasks.Run()
		assemble(lfilename)

	}
	hash, _ := security.FileHash256(lfilename)
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
	descriptor *model.FileDescriptor
	target     string
	cp         IMessageHandler
	service    *FileManager
}

func NewFetchPartTask(descriptor *model.FileDescriptor, target string, cp IMessageHandler, service *FileManager) *FetchPartTask {
	fpt := &FetchPartTask{}
	fpt.descriptor = descriptor
	fpt.target = target
	fpt.cp = cp
	fpt.service = service
	return fpt
}

func (task *FetchPartTask) Run() {
	data := task.cp.Request(task.descriptor, task.service.PeerServiceID()).(*model.FileData)
	file, _ := os.Create(task.target + ".part-" + getPartString(task.descriptor.Part()))
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
