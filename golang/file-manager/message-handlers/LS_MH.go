package message_handlers

import (
	. "github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type LS_MH struct {
	fs *FileManagerService
}

func NewRlsMH(service IService) *LS_MH {
	lf := &LS_MH{}
	lf.fs = service.(*FileManagerService)
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
	fr := &model.FileRequest{}
	fr.UnMarshal(message.Data())
	fd := model.NewFileDescriptor(fr.Path(), fr.Dept())
	if fd == nil {
		fd = &model.FileDescriptor{}
	}
	data := fd.Marshal()
	msgHandler.fs.ServiceManager().Reply(message, data)
}

func (msgHandler *LS_MH) Request(data interface{}, dest *ServiceID) interface{} {
	req := data.(*model.FileRequest)
	response, e := msgHandler.fs.ServiceManager().Request(msgHandler.Topic(), msgHandler.fs, dest, req.Marshal(), false)
	if e != nil {
		return nil
	}
	fd := model.UnmarshalFileDescriptor(response)
	return fd
}
