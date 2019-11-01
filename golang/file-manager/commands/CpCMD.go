package commands

import (
	"bytes"
	. "github.com/saichler/console/golang/console/commands"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	"io/ioutil"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
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
		job := &Job{}
		job.mtx = sync.NewCond(&sync.Mutex{})
		job.parts = parts
		job.mtx.L.Lock()
		for i := 0; i < parts; i++ {
			descriptor := fd.Clone()
			descriptor.SetPart(i)
			go c.FetchPart(descriptor, lfilename, job)
		}
		for job.parts != 0 {
			job.mtx.Wait()
		}
		job.mtx.L.Unlock()
		assemble(lfilename)

	}
	return "Multipart File with:" + strconv.Itoa(parts) + " and the size of:" + strconv.Itoa(int(fd.Size())), nil
}

func (c *CpCMD) FetchPart(descriptor *model.FileDescriptor, target string, job *Job) {
	data := c.cp.Request(descriptor, c.service.PeerServiceID()).(*model.FileData)
	file, _ := os.Create(target + ".part-" + getPartString(descriptor.Part()))
	file.Write(data.Data())
	file.Close()
	job.mtx.L.Lock()
	job.parts--
	if job.parts == 0 {
		job.mtx.Broadcast()
	}
	job.mtx.L.Unlock()
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

type Job struct {
	mtx   *sync.Cond
	parts int
}
