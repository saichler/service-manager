package model

import (
	"github.com/saichler/messaging/golang/net/protocol"
	"github.com/saichler/security"
	utils "github.com/saichler/utils/golang"
)

type Inventory struct {
	SID      *protocol.ServiceID
	Services []*protocol.ServiceID
}

func (inv *Inventory) Marshal(oldHash string) ([]byte, string) {
	bs := utils.NewByteSlice()
	inv.SID.Marshal(bs)
	bs.AddInt(len(inv.Services))
	for _, sid := range inv.Services {
		sid.Marshal(bs)
	}
	data := bs.Data()
	hash := security.Hash256(data)
	if oldHash != hash {
		return data, hash
	}

	bs = utils.NewByteSlice()
	inv.SID.Marshal(bs)
	bs.AddInt(0)
	return bs.Data(), oldHash
}

func (inv *Inventory) UnMarshal(data []byte) {
	bs := utils.NewByteSliceWithData(data, 0)
	inv.SID = &protocol.ServiceID{}
	inv.SID.Unmarshal(bs)
	size := bs.GetInt()
	inv.Services = make([]*protocol.ServiceID, size)
	for i := 0; i < size; i++ {
		inv.Services[i] = &protocol.ServiceID{}
		inv.Services[i].Unmarshal(bs)
	}
}
