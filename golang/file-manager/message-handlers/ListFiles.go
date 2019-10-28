package message_handlers

import (
	"fmt"
	. "github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type ListFiles struct {
	fs *FileManager
}

func NewListFiles(service IService) *ListFiles {
	lf := &ListFiles{}
	lf.fs = service.(*FileManager)
	return lf
}

func (m *ListFiles) Init() {
}

func (m *ListFiles) Topic() string {
	return "ListFiles"
}

func (m *ListFiles) Message(destination *ServiceID, data []byte, isReply bool) *Message {
	msg := m.fs.ServiceManager().NewMessage(m.Topic(), m.fs, destination, data, isReply)
	return msg
}

func (m *ListFiles) Handle(message *Message) {
	fmt.Println("Received list files")
}
