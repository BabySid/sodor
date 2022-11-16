package routine

import (
	"github.com/BabySid/gobase"
	"sodor/thomas/fat_ctrl"
	"sync"
)

type routine struct {
	scheduler *gobase.Scheduler
}

var (
	once      sync.Once
	singleton *routine
)

func GetInstance() *routine {
	once.Do(func() {
		singleton = &routine{
			scheduler: gobase.NewScheduler(),
		}
	})
	return singleton
}

func (r *routine) initJobs() error {
	err := r.scheduler.AddJob("handShakeWithFatCtrl", "*/20 * * * * *", fat_ctrl.GetInstance())
	if err != nil {
		return err
	}

	return nil
}

func (r *routine) Start() error {
	r.scheduler.Start()
	err := r.initJobs()
	return err
}
