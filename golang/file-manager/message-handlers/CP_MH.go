package message_handlers

import (
	. "github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type CP_MH struct {
	fs *FileManagerService
}

func NewCpMH(service IService) *CP_MH {
	lf := &CP_MH{}
	lf.fs = service.(*FileManagerService)
	return lf
}

func (msgHandler *CP_MH) Init() {
}

func (msgHandler *CP_MH) Topic() string {
	return "cp"
}

func (msgHandler *CP_MH) Message(destination *ServiceID, data []byte, isReply bool) *Message {
	msg := msgHandler.fs.ServiceManager().NewMessage(msgHandler.Topic(), msgHandler.fs, destination, data, isReply)
	return msg
}

var total = 0

func (msgHandler *CP_MH) Handle(message *Message) {
	fileData := &model.FileData{}
	fileData.Unmarshal(message.Data())
	fileData.LoadData()
	total += len(fileData.Data())
	msgHandler.fs.ServiceManager().Reply(message, fileData.Marshal())
}

func (msgHandler *CP_MH) Request(data interface{}, dest *ServiceID) interface{} {
	fileData := data.(*model.FileData)
	response, e := msgHandler.fs.ServiceManager().Request(msgHandler.Topic(), msgHandler.fs, dest, fileData.Marshal(), false)
	if e != nil {
		return nil
	}
	fileData.Unmarshal(response)
	return fileData
}
