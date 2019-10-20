package service_manager

import (
	"github.com/saichler/console/golang/console/commands"
	. "github.com/saichler/service-manager/golang/common"
	"strconv"
	"sync"
)

type ServiceContainer struct {
	topic            string
	nextServiceID    uint16
	serviceInstances map[uint16]IService
	mtx              *sync.Mutex
}

func NewServiceContainer(topic string) *ServiceContainer {
	sc := &ServiceContainer{}
	sc.topic = topic
	sc.serviceInstances = make(map[uint16]IService)
	sc.mtx = &sync.Mutex{}
	return sc
}

func (sc *ServiceContainer) Topic() string {
	return sc.topic
}

func (sc *ServiceContainer) AddService(service IService, serviceManager *ServiceManager) error {
	sc.mtx.Lock()
	defer sc.mtx.Unlock()
	sc.nextServiceID++
	sci := commands.NewConsoleID(service.Topic()+"-"+strconv.Itoa(int(sc.nextServiceID)), serviceManager.cid)
	service.Init(serviceManager, sc.nextServiceID, sci)
	sc.serviceInstances[sc.nextServiceID] = service
	service.Start()
	return nil
}

func (sc *ServiceContainer) Services() []IService {
	list := make([]IService, 0)
	for _, service := range sc.serviceInstances {
		list = append(list, service)
	}
	return list
}

func (sc *ServiceContainer) Service(id uint16) IService {
	return sc.serviceInstances[id]
}
