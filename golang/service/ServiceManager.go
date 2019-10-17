package service

import (
	. "github.com/saichler/messaging/golang/net/netnode"
	. "github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/service/console"
	. "github.com/saichler/utils/golang"
)

type ServiceManager struct {
	node       *NetworkNode
	console    *console.Console
	containers map[string]*ServiceContainer
}

func NewServiceManager() (*ServiceManager, error) {
	sm := &ServiceManager{}
	sm.containers = make(map[string]*ServiceContainer)
	node, e := NewNetworkNode(sm)
	if e != nil {
		Error("Failed to create a network node:", e)
		return nil, e
	}
	sm.node = node
	sm.console = console.NewConsole(node.Port()-10000)
	return sm, nil
}

func (sm *ServiceManager) HandleMessage(message *Message) {

}

func (sm *ServiceManager) HandleUnreachable(message *Message) {

}
