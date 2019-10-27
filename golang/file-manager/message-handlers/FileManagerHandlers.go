package message_handlers

import . "github.com/saichler/service-manager/golang/service-manager"

type FileManagerHandlers struct{}

func (mh *FileManagerHandlers) Handlers(service IService) []IMessageHandler {
	result := make([]IMessageHandler, 0)
	return result
}
