package service_manager

import (
	"fmt"
	. "github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/common"
	. "github.com/saichler/utils/golang"
)

func (sm *ServiceManager) processMessages() {
	for sm.node.Running() {
		message, ok := sm.inbox.Pop().(*Message)
		if ok {
			sm.processMessage(message)
		}
	}
}

func (sm *ServiceManager) processMessage(message *Message) {
	if message.Destination().Publish() {
		fmt.Println("Got Publish Message with topic:" + message.Topic())
	} else if message.Destination().Unreachable() {
		fmt.Println("Got Unreachable Message with topic:" + message.Topic())
	} else {
		service, e := sm.getServiceForMessage(message)
		if e != nil {
			return
		}
		service.Handle(message)
	}
}

func (sm *ServiceManager) getServiceForMessage(message *Message) (IService, error) {
	container, ok := sm.containers[message.Destination().Topic()]
	if ok {
		return container.Service(message.Destination().ID()), nil
	} else {
		return nil, Error("Unknown Service Type" + message.Destination().Topic())
	}
	return nil, nil
}
