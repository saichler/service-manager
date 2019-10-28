package message_handlers

import (
	service2 "github.com/saichler/service-manager/golang/management/service"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type ManagementHandlers struct {
	service  IService
	handlers map[string]IMessageHandler
}

func (mh *ManagementHandlers) Init(service IService) {
	mh.service = service
	mh.handlers = make(map[string]IMessageHandler)
	mh.addHandler(NewPingMH(service.(*service2.ManagementService)))
}

func (mh *ManagementHandlers) Handlers() []IMessageHandler {
	result := make([]IMessageHandler, 0)
	for _, h := range mh.handlers {
		result = append(result, h)
	}
	return result
}

func (mh *ManagementHandlers) Handler(topic string) IMessageHandler {
	return mh.handlers[topic]
}

func (mh *ManagementHandlers) addHandler(handler IMessageHandler) {
	mh.handlers[handler.Topic()] = handler
}
