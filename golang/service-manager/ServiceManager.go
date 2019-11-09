package service_manager

import (
	. "github.com/saichler/console/golang/console"
	"github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/messaging/golang/net/netnode"
	. "github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-management-service/golang/management-service/model"

	//	commands2 "github.com/saichler/service-manager/golang/management/commands"
	. "github.com/saichler/utils/golang"
	"strconv"
	"sync"
	"time"
)

type ServiceManager struct {
	cid             *commands.ConsoleId
	node            *NetworkNode
	console         *Console
	containers      *MapList
	shutdownMtx     *sync.Cond
	inbox           *PriorityQueue
	msgScheduler    *MessageScheduler
	active          bool
	serviceNetwork  *model.ServiceNetwork
	pendingRequests *Map
}

func NewServiceManager() (*ServiceManager, error) {
	sm := &ServiceManager{}
	sm.serviceNetwork = model.NewServiceNetwork()
	sm.active = true
	sm.containers = NewMapList()
	sm.pendingRequests = NewMap()
	sm.inbox = NewPriorityQueue()
	sm.msgScheduler = newMessageScheduler()
	node, e := NewNetworkNode(sm)
	if e != nil {
		Error("Failed to create a network node:", e)
		return nil, e
	}
	sm.node = node
	sm.cid = commands.NewConsoleID(GetIpAsString(node.NetworkID().Host())+":"+strconv.Itoa(int(node.NetworkID().Port())), nil)
	sm.shutdownMtx = sync.NewCond(&sync.Mutex{})

	sm.console, _ = NewConsole("127.0.0.1", int(node.Port())-10000, sm.cid)
	Info("Console bind to 127.0.0.1:" + strconv.Itoa(int(node.Port())-10000))
	sm.console.Start(false)

	go sm.processMessages()
	go sm.runScheduler()

	return sm, nil
}

func (sm *ServiceManager) HandleMessage(message *Message) {
	sm.inbox.Push(message, message.Priority())
}

func (sm *ServiceManager) HandleUnreachable(message *Message) {

}

func (sm *ServiceManager) WaitForShutdown() {
	sm.shutdownMtx.L.Lock()
	sm.shutdownMtx.Wait()
}

func (sm *ServiceManager) Shutdown() {
	sm.active = false
	go sm.shutdown()
}

func (sm *ServiceManager) shutdown() {
	Info("Initiating shutdown")
	time.Sleep(time.Second * 2)
	sm.node.Shutdown()
	time.Sleep(time.Second * 2)
	sm.shutdownMtx.Broadcast()
}

func (sm *ServiceManager) ConsoleId() *commands.ConsoleId {
	return sm.cid
}

func (sm *ServiceManager) Console() *Console {
	return sm.console
}

func (sm *ServiceManager) Services() []IService {
	containers := sm.containers.List()
	list := make([]IService, 0)
	for _, sc := range containers {
		container := sc.(*ServiceContainer)
		scs := container.Services()
		list = append(list, scs...)
	}
	return list
}

func (sm *ServiceManager) NetworkID() *NetworkID {
	return sm.node.NetworkID()
}

func (sm *ServiceManager) Publish(topic string, service IService, data []byte) error {
	dest := NewPublishServiceID(topic)
	message := sm.node.NewMessage(service.ServiceID(), dest, service.ServiceID(), topic, 0, data, false)
	e := sm.node.SendMessage(message)
	if e != nil {
		return Error("Failed to send message")
	}
	return nil
}

func (sm *ServiceManager) NewMessage(msgTopic string, source IService, destination *ServiceID, data []byte, isReply bool) *Message {
	return sm.node.NewMessage(source.ServiceID(), destination, source.ServiceID(), msgTopic, 0, data, isReply)
}

func (sm *ServiceManager) Reply(message *Message, data []byte) error {
	msg := NewMessage(message.Destination(), message.Source(), message.OriginalSource(), message.MessageID(), message.Topic(), message.Priority(), data, true)
	e := sm.node.SendMessage(msg)
	if e != nil {
		return Error("Failed to send message")
	}
	return nil
}

func (sm *ServiceManager) Send(topic string, source IService, destination *ServiceID, data []byte, isReply bool) error {
	message := sm.NewMessage(topic, source, destination, data, isReply)
	e := sm.node.SendMessage(message)
	if e != nil {
		return Error("Failed to send message")
	}
	return nil
}

func (sm *ServiceManager) Uplink(ip string) {
	sm.node.Uplink(ip)
}

func (sm *ServiceManager) ServiceNetwork() *model.ServiceNetwork {
	return sm.serviceNetwork
}
