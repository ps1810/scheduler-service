package local_cron

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type LocalCron struct {
	Cron   *cron.Cron
	JobMap map[uint]cron.EntryID
}

var m sync.Mutex
var localCron LocalCron

func InitCron() {
	m.Lock()
	defer m.Unlock()

	if localCron.Cron == nil {
		localCron.Cron = cron.New()
		localCron.JobMap = make(map[uint]cron.EntryID)
	}
}

func GetCron() *LocalCron {
	if localCron.Cron == nil {
		InitCron()
		StartCron()
	}
	return &localCron
}

func StartCron() {
	localCron.Cron.Start()
}

func StopCron() {
	localCron.Cron.Stop()
}
