package service_manager

import (
	"github.com/saichler/console/golang/console/commands"
	"github.com/saichler/messaging/golang/net/protocol"
	"strconv"
	"sync"
)

type ServiceContainer struct {
	topic            string
	nextServiceID    uint16
	serviceInstances map[uint16]*ServiceContainerEntry
	mtx              *sync.Mutex
}

func NewServiceContainer(topic string) *ServiceContainer {
	sc := &ServiceContainer{}
	sc.topic = topic
	sc.serviceInstances = make(map[uint16]*ServiceContainerEntry)
	sc.mtx = &sync.Mutex{}
	return sc
}

func (sc *ServiceContainer) Topic() string {
	return sc.topic
}

func (sc *ServiceContainer) AddService(service IService, handlers IServiceMessageHandlers, serviceManager *ServiceManager) error {
	sc.mtx.Lock()
	defer sc.mtx.Unlock()
	sc.nextServiceID++
	sci := commands.NewConsoleID(service.Topic()+"-"+strconv.Itoa(int(sc.nextServiceID)), serviceManager.cid)
	sid := protocol.NewServiceID(serviceManager.NetworkID(), service.Topic(), sc.nextServiceID)
	service.Init(serviceManager, sc.nextServiceID, sid, sci)
	sc.serviceInstances[sc.nextServiceID] = newServiceContainerEntry(service)
	for _, handler := range handlers.Handlers(service) {
		sc.serviceInstances[sc.nextServiceID].RegisterMessageHandler(handler)
		handler.Init()
	}
	sc.serviceInstances[sc.nextServiceID].start()
	return nil
}

func (sc *ServiceContainer) Services() []IService {
	sc.mtx.Lock()
	defer sc.mtx.Unlock()
	list := make([]IService, 0)
	for _, serviceEntry := range sc.serviceInstances {
		list = append(list, serviceEntry.service)
	}
	return list
}

func (sc *ServiceContainer) ServiceContainerEntry(id uint16) *ServiceContainerEntry {
	sc.mtx.Lock()
	defer sc.mtx.Unlock()
	return sc.serviceInstances[id]
}

func (sc *ServiceContainer) publish(message *protocol.Message) {
	sc.mtx.Lock()
	defer sc.mtx.Unlock()
	for _, se := range sc.serviceInstances {
		se.inbox.Push(message, message.Priority())
	}
}
