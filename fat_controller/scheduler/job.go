package scheduler

import (
	"github.com/BabySid/proto/sodor"
	"github.com/robfig/cron/v3"
)

type jobContext struct {
	job    *sodor.Job
	cronID cron.EntryID
}

func (j *jobContext) run() error {
	return nil
}
