package service

import (
	. "github.com/saichler/console/golang/console"
	"github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/messaging/golang/net/netnode"
	. "github.com/saichler/messaging/golang/net/protocol"
	commands2 "github.com/saichler/service-manager/golang/commands"
	. "github.com/saichler/utils/golang"
	"strconv"
	"sync"
)

type ServiceManager struct {
	cid         *commands.ConsoleId
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
	sm.cid = commands.NewConsoleID(GetIpAsString(node.NetworkID().Host())+":"+strconv.Itoa(int(node.NetworkID().Port())), nil)
	sm.shutdownMtx = sync.NewCond(&sync.Mutex{})
	sm.console, _ = NewConsole("127.0.0.1", node.Port()-10000, sm.cid)
	sm.console.RegisterCommand(commands2.NewShutdown(sm))
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

func (sm *ServiceManager) Shutdown() {
	sm.shutdownMtx.Broadcast()
}

func (sm *ServiceManager) ConsoleId() *commands.ConsoleId {
	return sm.cid
}
