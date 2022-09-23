package scheduler

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sync"
)

type scheduler struct {
	routine     *cron.Cron
	builtInJobs map[string]cron.EntryID // name => entry
}

var (
	once      sync.Once
	singleton *scheduler
)

func GetInstance() *scheduler {
	once.Do(func() {
		singleton = &scheduler{}
		err := singleton.initOnce()
		if err != nil {
			log.Fatalf("scheduler init failed. err=%s", err)
		}
	})
	return singleton
}

func (s *scheduler) Start() error {
	s.routine.Start()
	return nil
}

func (s *scheduler) initOnce() error {
	s.routine = cron.New(cron.WithSeconds(), cron.WithLogger(cron.PrintfLogger(&cronLog{})))
	return nil
}

type cronLog struct{}

func (l *cronLog) Printf(msg string, v ...interface{}) {
	log.Infof(msg, v...)
}
