package message_handlers

import (
	"github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/service-manager/golang/management/model"
	. "github.com/saichler/service-manager/golang/management/service"
	. "github.com/saichler/service-manager/golang/service-manager"
	utils "github.com/saichler/utils/golang"
)

type PingMH struct {
	ms   *ManagementService
	ping *protocol.Message
	hash string
}

func NewPingMH(service IService) *PingMH {
	ping := &PingMH{}
	ping.ms = service.(*ManagementService)
	return ping
}

func (m *PingMH) Init() {
	m.ms.ServiceManager().ScheduleMessage(m, 10, 0)
}

func (m *PingMH) Topic() string {
	return "Ping"
}

func (m *PingMH) Message(destination *protocol.ServiceID, data []byte, isReply bool) *protocol.Message {
	dest := protocol.NewServiceID(protocol.NetConfig.PublishID(), m.ms.Topic(), m.ms.ID())
	return m.ms.ServiceManager().NewMessage(m.Topic(), m.ms, dest, m.inventory(), false)
}

func (m *PingMH) Handle(message *protocol.Message) {
	inv := &model.Inventory{}
	inv.UnMarshal(message.Data())
	utils.Info("Reveived Inventory From:", message.Source().String(), " with:")
	for _, s := range inv.Services {
		utils.Info("  ", s.String())
	}
	m.ms.ServiceManager().ServiceNetwork().UpdateInventory(inv)
}

func (m *PingMH) inventory() []byte {
	inv := &model.Inventory{}
	inv.SID = m.ms.ServiceID()
	services := m.ms.ServiceManager().Services()
	inv.Services = make([]*protocol.ServiceID, 0)
	for _, service := range services {
		inv.Services = append(inv.Services, service.ServiceID())
	}
	data, hash := inv.Marshal(m.hash)
	m.hash = hash
	return data
}

func (m *PingMH) Request(data interface{}, destination *protocol.ServiceID) interface{}{
	return nil
}
