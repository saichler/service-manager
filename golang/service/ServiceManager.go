package service

import (
	. "github.com/saichler/console/golang/console"
	. "github.com/saichler/messaging/golang/net/netnode"
	. "github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/utils/golang"
	"strconv"
	"sync"
)

type ServiceManager struct {
	node        *NetworkNode
	console     *Console
	containers  map[string]*ServiceContainer
	shutdownMtx *sync.Cond
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
	sm.shutdownMtx = sync.NewCond(&sync.Mutex{})
	sm.console, _ = NewConsole("127.0.0.1", node.Port()-10000, sm)
	Info("Console bind to 127.0.0.1:" + strconv.Itoa(node.Port()-10000))
	sm.console.Start(false)
	return sm, nil
}

func (sm *ServiceManager) HandleMessage(message *Message) {

}

func (sm *ServiceManager) HandleUnreachable(message *Message) {

}

func (sm *ServiceManager) WaitForShutdown() {
	sm.shutdownMtx.L.Lock()
	sm.shutdownMtx.Wait()
}
