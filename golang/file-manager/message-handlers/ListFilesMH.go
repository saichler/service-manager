package message_handlers

import (
	. "github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type ListFilesMH struct {
	fs *FileManager
}

func NewListFilesMH(service IService) *ListFilesMH {
	lf := &ListFilesMH{}
	lf.fs = service.(*FileManager)
	return lf
}

func (m *ListFilesMH) Init() {
}

func (m *ListFilesMH) Topic() string {
	return "ListFiles"
}

func (m *ListFilesMH) Message(destination *ServiceID, data []byte, isReply bool) *Message {
	msg := m.fs.ServiceManager().NewMessage(m.Topic(), m.fs, destination, data, isReply)
	return msg
}

func (m *ListFilesMH) Handle(message *Message) {
	if message.IsReply() {
		panic("rrrrrr")
	}
	dir := string(message.Data())
	fd, _ := model.Create(dir)
	data := fd.Marshal()
	m.fs.ServiceManager().Reply(message, data)
}

func (m *ListFilesMH) Request(data interface{}, dest *ServiceID) interface{} {
	response, e := m.fs.ServiceManager().Request(m.Topic(), m.fs, dest, []byte(data.(string)), false)
	if e != nil {
		return nil
	}
	fd := &model.FileDescriptor{}
	fd.Unmarshal(response)
	return fd
}
