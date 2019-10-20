package service

import (
	"github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/management/model"
)

type InventoryMessage struct {
	ms   *ManagementService
	ping *protocol.Message
}

func NewInventoryMessage(ms *ManagementService) *InventoryMessage {
	inv := &InventoryMessage{}
	inv.ms = ms
	return inv
}

func (m *InventoryMessage) Send() {
	m.ms.ServiceManager().Publish("Inventory", m.ms, m.inventory())
}

func (m *InventoryMessage) Message() *protocol.Message {
	dest := protocol.NewServiceID(protocol.NetConfig.PublishID(), "Ping", m.ms.id)
	return m.ms.sm.NewMessage("Ping", m.ms, dest, m.inventory())
}

func (m *InventoryMessage) Handle(message *protocol.Message) {

}

func (m *InventoryMessage) inventory() []byte {
	inv := &model.Inventory{}
	inv.SID = m.ms.serviceID
	services := m.ms.ServiceManager().Services()
	inv.Services = make([]*protocol.ServiceID, 0)
	for _, service := range services {
		inv.Services = append(inv.Services, service.ServiceID())
	}
	return inv.Marshal()
}
