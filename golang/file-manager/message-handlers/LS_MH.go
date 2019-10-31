package message_handlers

import (
	. "github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type LS_MH struct {
	fs *FileManager
}

func NewRlsMH(service IService) *LS_MH {
	lf := &LS_MH{}
	lf.fs = service.(*FileManager)
	return lf
}

func (msgHandler *LS_MH) Init() {
}

func (msgHandler *LS_MH) Topic() string {
	return "ls"
}

func (msgHandler *LS_MH) Message(destination *ServiceID, data []byte, isReply bool) *Message {
	msg := msgHandler.fs.ServiceManager().NewMessage(msgHandler.Topic(), msgHandler.fs, destination, data, isReply)
	return msg
}

func (msgHandler *LS_MH) Handle(message *Message) {
	dir := string(message.Data())
	fd, _ := model.Create(dir, 1, 0)
	if fd == nil {
		fd = &model.FileDescriptor{}
	}
	data := fd.Marshal()
	msgHandler.fs.ServiceManager().Reply(message, data)
}

func (msgHandler *LS_MH) Request(data interface{}, dest *ServiceID) interface{} {
	response, e := msgHandler.fs.ServiceManager().Request(msgHandler.Topic(), msgHandler.fs, dest, []byte(data.(string)), false)
	if e != nil {
		return nil
	}
	fd := &model.FileDescriptor{}
	fd.Unmarshal(response)
	return fd
}
