package service_manager

import (
	"github.com/saichler/service-manager/golang/common"
	"sync"
	"time"
)

type MessageScheduler struct {
	schedules []*MessageSchedulerEntry
	mtx       *sync.Mutex
}

type MessageSchedulerEntry struct {
	handler  common.IMessageHandler
	interval int64
	last     int64
}

func newMessageScheduler() *MessageScheduler {
	ms := &MessageScheduler{}
	ms.mtx = &sync.Mutex{}
	ms.schedules = make([]*MessageSchedulerEntry, 0)
	return ms
}

func (sm *ServiceManager) ScheduleMessage(handler common.IMessageHandler, interval, initial int64) {
	entry := &MessageSchedulerEntry{}
	entry.interval = interval
	entry.handler = handler
	entry.last = time.Now().Unix() - interval + initial
	sm.msgScheduler.mtx.Lock()
	defer sm.msgScheduler.mtx.Unlock()
	sm.msgScheduler.schedules = append(sm.msgScheduler.schedules, entry)
}

func (sm *ServiceManager) runScheduler() {
	for sm.active {
		sm.msgScheduler.mtx.Lock()
		for _, entry := range sm.msgScheduler.schedules {
			if time.Now().Unix() > entry.last+entry.interval {
				entry.last = time.Now().Unix()
				sm.node.SendMessage(entry.handler.Message())
			}
		}
		time.Sleep(time.Second)
		sm.msgScheduler.mtx.Unlock()
	}
}
