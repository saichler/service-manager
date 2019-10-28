package service_manager

import (
	"fmt"
	. "github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/utils/golang"
	"sync"
)

func (sm *ServiceManager) processMessages() {
	for sm.active {
		message, ok := sm.inbox.Pop().(*Message)
		if ok {
			sm.processMessage(message)
		}
	}
}

func (sm *ServiceManager) processMessage(message *Message) {
	condInt := sm.pendingRequests.Get(message.MessageID())
	if condInt != nil && message.IsReply() {
		cond := condInt.(*sync.Cond)
		cond.L.Lock()
		sm.pendingRequests.Put(message.MessageID(), message)
		cond.L.Unlock()
		cond.Broadcast()
		return
	}
	if message.Destination().Publish() {
		sm.handlePublish(message)
	} else if message.Destination().Unreachable() {
		fmt.Println("Got Unreachable Message with topic:" + message.Topic())
	} else {
		se, e := sm.getServiceEntryForMessage(message)
		if e != nil {
			return
		}
		se.inbox.Push(message, message.Priority())
	}
}

func (sm *ServiceManager) handlePublish(message *Message) {
	container, ok := sm.containers.Get(message.Destination().Topic()).(*ServiceContainer)
	if ok {
		container.publish(message)
	}
}

func (sm *ServiceManager) getServiceEntryForMessage(message *Message) (*ServiceContainerEntry, error) {
	container, ok := sm.containers.Get(message.Destination().Topic()).(*ServiceContainer)
	if ok {
		return container.ServiceContainerEntry(message.Destination().ID()), nil
	} else {
		return nil, Error("Unknown Service Type:" + message.Destination().Topic())
	}
	return nil, nil
}

func (sm *ServiceManager) Request(topic string, source IService, destination *ServiceID, data []byte, isReply bool) ([]byte, error) {
	message := sm.NewMessage(topic, source, destination, data, isReply)
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()
	sm.pendingRequests.Put(message.MessageID(), cond)
	e := sm.node.SendMessage(message)
	if e != nil {
		sm.pendingRequests.Del(message.MessageID())
		cond.L.Unlock()
		return nil, e
	}
	cond.Wait()
	msg := sm.pendingRequests.Del(message.MessageID()).(*Message)
	return msg.Data(), nil
}
