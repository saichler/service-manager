package message_handlers

import (
	service2 "github.com/saichler/service-manager/golang/management/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type ManagementHandlers struct{}

func (mh *ManagementHandlers) Handlers(service IService) []IMessageHandler {
	result := make([]IMessageHandler, 0)
	result = append(result, NewPingMH(service.(*service2.ManagementService)))
	return result
}
