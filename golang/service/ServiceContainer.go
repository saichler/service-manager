package service

import (
	. "github.com/saichler/service-manager/golang/common"
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
	service.Init(serviceManager, sc.nextServiceID)
	sc.serviceInstances[sc.nextServiceID] = service
	service.Start()
	return nil
}
