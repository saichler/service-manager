package message_handlers

import (
	. "github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type RlsMH struct {
	fs *FileManager
}

func NewRlsMH(service IService) *RlsMH {
	lf := &RlsMH{}
	lf.fs = service.(*FileManager)
	return lf
}

func (msgHandler *RlsMH) Init() {
}

func (msgHandler *RlsMH) Topic() string {
	return "RLS"
}

func (msgHandler *RlsMH) Message(destination *ServiceID, data []byte, isReply bool) *Message {
	msg := msgHandler.fs.ServiceManager().NewMessage(msgHandler.Topic(), msgHandler.fs, destination, data, isReply)
	return msg
}

func (msgHandler *RlsMH) Handle(message *Message) {
	dir := string(message.Data())
	fd, _ := model.Create(dir, 1, 0)
	data := fd.Marshal()
	msgHandler.fs.ServiceManager().Reply(message, data)
}

func (msgHandler *RlsMH) Request(data interface{}, dest *ServiceID) interface{} {
	response, e := msgHandler.fs.ServiceManager().Request(msgHandler.Topic(), msgHandler.fs, dest, []byte(data.(string)), false)
	if e != nil {
		return nil
	}
	fd := &model.FileDescriptor{}
	fd.Unmarshal(response)
	return fd
}
