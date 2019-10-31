package service

import (
	"github.com/saichler/console/golang/console/commands"
	"github.com/saichler/messaging/golang/net/protocol"
	. "github.com/saichler/service-manager/golang/service-manager"
)

type FileManager struct {
	id       uint16
	cid      *commands.ConsoleId
	sm       *ServiceManager
	sid      *protocol.ServiceID
	peerSID  *protocol.ServiceID
	peerDir  string
	localDir string
}

func (fm *FileManager) Topic() string {
	return "File Manager"
}

func (fm *FileManager) ID() uint16 {
	return fm.id
}

func (fm *FileManager) ConsoleId() *commands.ConsoleId {
	return fm.cid
}

func (fm *FileManager) ServiceManager() *ServiceManager {
	return fm.sm
}

func (fm *FileManager) ServiceID() *protocol.ServiceID {
	return fm.sid
}

func (fm *FileManager) Init(sm *ServiceManager, id uint16, sid *protocol.ServiceID, cid *commands.ConsoleId) {
	fm.sm = sm
	fm.id = id
	fm.cid = cid
	fm.sid = sid
	fm.peerSID = sid
	fm.peerDir = "/tmp"
	fm.localDir = "/tmp"
}

func (fm *FileManager) PeerServiceID() *protocol.ServiceID {
	return fm.peerSID
}

func (fm *FileManager) SetPeerServiceID(sid *protocol.ServiceID) {
	fm.peerSID = sid
}

func (fm *FileManager) PeerDir() string {
	return fm.peerDir
}

func (fm *FileManager) SetPeerDir(dir string) {
	fm.peerDir = dir
}

func (fm *FileManager) LocalDir() string {
	return fm.localDir
}
