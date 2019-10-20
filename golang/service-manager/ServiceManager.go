package service_manager

import (
	. "github.com/saichler/console/golang/console"
	"github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/messaging/golang/net/netnode"
	. "github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/common"
	commands2 "github.com/saichler/service-manager/golang/management/commands"
	"github.com/saichler/service-manager/golang/management/service"
	. "github.com/saichler/utils/golang"
	"strconv"
	"sync"
	"time"
)

type ServiceManager struct {
	cid          *commands.ConsoleId
	node         *NetworkNode
	console      *Console
	containers   map[string]*ServiceContainer
	shutdownMtx  *sync.Cond
	mtx          *sync.Mutex
	inbox        *PriorityQueue
	msgScheduler *MessageScheduler
}

func NewServiceManager() (*ServiceManager, error) {
	sm := &ServiceManager{}
	sm.mtx = &sync.Mutex{}
	sm.inbox = NewPriorityQueue()
	sm.containers = make(map[string]*ServiceContainer)
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

	sm.console.RegisterCommand(commands2.NewLS(sm))
	sm.console.RegisterCommand(commands2.NewCD(sm))

	go sm.processMessages()
	go sm.runScheduler()

	sm.AddService(&service.ManagementService{})

	return sm, nil
}

func (sm *ServiceManager) AddService(service common.IService) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	container, ok := sm.containers[service.Topic()]
	if !ok {
		container = NewServiceContainer(service.Topic())
		sm.containers[service.Topic()] = container
	}
	container.AddService(service, sm)
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

func (sm *ServiceManager) Services() []common.IService {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	list := make([]common.IService, 0)
	for _, sc := range sm.containers {
		scs := sc.Services()
		list = append(list, scs...)
	}
	return list
}

func (sm *ServiceManager) NetworkID() *NetworkID {
	return sm.node.NetworkID()
}

func (sm *ServiceManager) Publish(topic string, service common.IService, data []byte) error {
	dest := NewPublishServiceID(topic)
	message := sm.node.NewMessage(service.ServiceID(), dest, service.ServiceID(), topic, 0, data)
	e := sm.node.SendMessage(message)
	if e != nil {
		return Error("Failed to send message")
	}
	return nil
}

func (sm *ServiceManager) NewMessage(topic string, source common.IService, destination *ServiceID, data []byte) *Message {
	return sm.node.NewMessage(source.ServiceID(), destination, source.ServiceID(), topic, 0, data)
}

func (sm *ServiceManager) Send(topic string, source common.IService, destination *ServiceID, data []byte) error {
	message := sm.NewMessage(topic, source, destination, data)
	e := sm.node.SendMessage(message)
	if e != nil {
		return Error("Failed to send message")
	}
	return nil
}

func (sm *ServiceManager) Uplink(ip string) {
	sm.node.Uplink(ip)
}
