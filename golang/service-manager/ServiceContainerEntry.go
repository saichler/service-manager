package service_manager

import (
	"github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/common"
	. "github.com/saichler/utils/golang"
)

type ServiceContainerEntry struct {
	service         IService
	inbox           *PriorityQueue
	messageHandlers map[string]IMessageHandler
	active          bool
}

func newServiceContainerEntry(service IService) *ServiceContainerEntry {
	sce := &ServiceContainerEntry{}
	sce.inbox = NewPriorityQueue()
	sce.service = service
	sce.messageHandlers = make(map[string]IMessageHandler)
	return sce
}

func (sce *ServiceContainerEntry) processMessages() {
	for sce.active {
		msg, ok := sce.inbox.Pop().(*protocol.Message)
		if ok {
			mh, ok := sce.messageHandlers[msg.Topic()]
			if ok {
				mh.Handle(msg)
			} else {
				Error("Cannot find message handler for topic:" + msg.Topic())
			}
		}
	}
}

func (sce *ServiceContainerEntry) start() {
	sce.active = true
	go sce.processMessages()
}

func (sce *ServiceContainerEntry) shutdown() {
	sce.active = false
	sce.inbox.Push("", 0)
}

func (sce *ServiceContainerEntry) RegisterMessageHandler(handler IMessageHandler) {
	sce.messageHandlers[handler.Topic()] = handler
}
