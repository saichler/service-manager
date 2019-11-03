package jobs

import (
	. "github.com/saichler/service-manager/golang/file-manager/model"
	. "github.com/saichler/service-manager/golang/file-manager/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	. "github.com/saichler/utils/golang"
)

type Copy struct {
	messageHandler IMessageHandler
	fileManager    *FileManagerService
	listener       JobListener
	source         *FileDescriptor
	dest           *FileDescriptor
}

func NewCopy(listener JobListener, messageHandler IMessageHandler, fileManager *FileManagerService, source *FileDescriptor, dest string) *Copy {
	cp := &Copy{}
	cp.fileManager = fileManager
	cp.messageHandler = messageHandler
	cp.listener = listener
	cp.source = source
	cp.dest = NewFileDescriptor(dest, 100)
	return cp
}

func (c *Copy) Run() {
	if c.source.Files()==nil {

	} else {

	}
}
