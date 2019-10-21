package service_manager

import (
	"errors"
	. "github.com/saichler/console/golang/console"
	"github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/messaging/golang/net/netnode"
	. "github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/common"
	commands2 "github.com/saichler/service-manager/golang/management/commands"
	"github.com/saichler/service-manager/golang/management/service"
	. "github.com/saichler/utils/golang"
	"plugin"
	"strconv"
	"sync"
	"time"
)

type ServiceManager struct {
	cid           *commands.ConsoleId
	node          *NetworkNode
	console       *Console
	containers    map[string]*ServiceContainer
	containersMtx *sync.Mutex
	shutdownMtx   *sync.Cond
	inbox         *PriorityQueue
	msgScheduler  *MessageScheduler
	active        bool
}

func NewServiceManager() (*ServiceManager, error) {
	sm := &ServiceManager{}
	sm.active = true
	sm.containers = make(map[string]*ServiceContainer)
	sm.containersMtx = &sync.Mutex{}
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

	sm.console.RegisterCommand(commands2.NewLS(sm))
	sm.console.RegisterCommand(commands2.NewCD(sm))

	go sm.processMessages()
	go sm.runScheduler()

	sm.AddService(&service.ManagementService{})
	sm.LoadService("/home/saichler/datasand/src/github.com/saichler/service-manager/golang/file-manager/plugin/FileService.so")
	return sm, nil
}

func (sm *ServiceManager) AddService(service common.IService) {
	sm.containersMtx.Lock()
	container, ok := sm.containers[service.Topic()]
	if !ok {
		container = NewServiceContainer(service.Topic())
		sm.containers[service.Topic()] = container
	}
	sm.containersMtx.Unlock()
	container.AddService(service, sm)
	sm.console.RegisterCommand(commands2.NewCD(service))
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

func (sm *ServiceManager) Services() []common.IService {
	containers := make([]*ServiceContainer, 0)
	sm.containersMtx.Lock()
	for _, container := range sm.containers {
		containers = append(containers, container)
	}
	sm.containersMtx.Unlock()

	list := make([]common.IService, 0)
	for _, sc := range containers {
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

func (sm *ServiceManager) NewMessage(msgTopic string, source common.IService, destination *ServiceID, data []byte) *Message {
	return sm.node.NewMessage(source.ServiceID(), destination, source.ServiceID(), msgTopic, 0, data)
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

func (sm *ServiceManager) LoadService(filename string) error {
	servicePlugin, e := plugin.Open(filename)
	if e != nil {
		Error("Unable to load serivce plugin:", e)
		return e
	}
	svr, e := servicePlugin.Lookup("Service")
	if e != nil {
		Error("Unable to find ServiceInstance in the library " + filename)
		Error("Make sure you have: var ServiceInstance Service = &<your service struct>{}")
		return e
	}

	servicePtr, ok := svr.(*common.IService)
	if !ok {
		msg := "Service is not of type IService, please check that it implements IService and that Service is a pointer."
		Error(msg)
		return errors.New(msg)
	}
	service := *servicePtr
	sm.AddService(service)
	return nil
}
