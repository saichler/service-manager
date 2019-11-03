package model

import (
	"github.com/saichler/messaging/golang/net/protocol"
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

func (sn *ServiceNetwork) UpdateInventory(inventory *Inventory) bool {
	sn.mtx.Lock()
	defer sn.mtx.Unlock()
	key := inventory.SID.String()
	if len(inventory.Services) > 0 {
		utils.Info("Updating inventory")
		sn.serviceNetwork[key] = inventory
	}
	if sn.serviceNetwork[key] == nil {
		return true
	}
	return false
}

func (sn *ServiceNetwork) GetPeers(id *protocol.ServiceID) []*protocol.ServiceID {
	sn.mtx.Lock()
	defer sn.mtx.Unlock()
	result := make([]*protocol.ServiceID, 0)
	for _, inv := range sn.serviceNetwork {
		for _, s := range inv.Services {
			if s.Topic() == id.Topic() && s.String() != id.String() {
				result = append(result, s)
			}
		}
	}
	return result
}
