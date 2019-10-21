package model

import "sync"

type ServiceNetwork struct {
	mtx            *sync.Mutex
	serviceNetwork map[string]*Inventory
}

func NewServiceNetwork() *ServiceNetwork {
	sn := &ServiceNetwork{}
	sn.serviceNetwork = make(map[string]*Inventory)
	sn.mtx = &sync.Mutex{}
	return sn
}

func (sn *ServiceNetwork) AddInventory(inventory *Inventory) {
	sn.mtx.Lock()
	defer sn.mtx.Unlock()
	key := inventory.SID.String()
	sn.serviceNetwork[key] = inventory
}
