package model

import (
	utils "github.com/saichler/utils/golang"
	"sync"
)

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

func (sn *ServiceNetwork) UpdateInventory(inventory *Inventory) {
	sn.mtx.Lock()
	defer sn.mtx.Unlock()
	key := inventory.SID.String()
	if len(inventory.Services) > 0 {
		utils.Info("Updating inventory")
		sn.serviceNetwork[key] = inventory
	}
}


