package service

import (
	"github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/management/model"
	utils "github.com/saichler/utils/golang"
)

type PingMH struct {
	ms   *ManagementService
	ping *protocol.Message
}

func NewPingMH(ms *ManagementService) *PingMH {
	ping := &PingMH{}
	ping.ms = ms
	return ping
}

func (m *PingMH) Send() {
	m.ms.ServiceManager().Publish("Inventory", m.ms, m.inventory())
}

func (m *PingMH) Topic() string {
	return "Ping"
}

func (m *PingMH) Message() *protocol.Message {
	dest := protocol.NewServiceID(protocol.NetConfig.PublishID(), m.ms.Topic(), m.ms.id)
	return m.ms.sm.NewMessage(m.Topic(), m.ms, dest, m.inventory())
}

func (m *PingMH) Handle(message *protocol.Message) {
	inv := &model.Inventory{}
	inv.UnMarshal(message.Data())
	utils.Info("Reveived Inventory From:" + message.Source().String())
	m.ms.serviceNetwork.AddInventory(inv)
}

func (m *PingMH) inventory() []byte {
	inv := &model.Inventory{}
	inv.SID = m.ms.serviceID
	services := m.ms.ServiceManager().Services()
	inv.Services = make([]*protocol.ServiceID, 0)
	for _, service := range services {
		inv.Services = append(inv.Services, service.ServiceID())
	}
	return inv.Marshal()
}
