package scheduler

import (
	"errors"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sodor/fat_controller/metastore"
	"sync"
)

type scheduler struct {
	routine *cron.Cron
	jobs    sync.Map // jID => entry
}

var (
	once      sync.Once
	singleton *scheduler

	NotRoutineJob = errors.New("not Routine Job")
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

	local, err := gobase.GetLocalIP()
	if err != nil {
		return err
	}
	states, err := metastore.GetInstance().SelectScheduler(local)
	if err != nil {
		return err
	}

	for _, stat := range states {
		var job sodor.Job
		job.Id = stat.JobID
		err = metastore.GetInstance().SelectJob(&job)
		if err != nil {
			return err
		}

		if job.ScheduleMode != sodor.ScheduleMode_SM_Crontab {
			log.Warnf("job schdulemode is not crontab: %v", job.ScheduleMode)
			continue
		}

		ctx := newJobContext()
		err = ctx.setJob(&job)
		if err != nil {
			log.Warnf("job init-set failed: %s", err)
			continue
		}
		cid, err := s.routine.AddJob(job.RoutineSpec.CtSpec, ctx)
		if err != nil {
			return err
		}
		ctx.cronID = cid

		s.jobs.Store(job.Id, ctx)
	}
	return nil
}

func (s *scheduler) initOnce() error {
	s.routine = cron.New(cron.WithParser(NewParser()), cron.WithLogger(cron.PrintfLogger(&cronLog{})))
	return nil
}

func (s *scheduler) AddJob(job *sodor.Job) error {
	if job.ScheduleMode != sodor.ScheduleMode_SM_Crontab {
		log.Warnf("job schdulemode is not crontab: %v", job.ScheduleMode)
		return nil
	}
	ctx := newJobContext()
	err := ctx.setJob(job)
	if err != nil {
		log.Warnf("job init-set failed: %s", err)
		return err
	}

	cid, err := s.routine.AddJob(job.RoutineSpec.CtSpec, ctx)
	if err != nil {
		return err
	}
	ctx.cronID = cid

	s.jobs.Store(job.Id, ctx)
	return nil
}

func (s *scheduler) Remove(job *sodor.Job) error {
	ctx, ok := s.jobs.LoadAndDelete(job.Id)
	if ok {
		jc := ctx.(*jobContext)
		s.routine.Remove(jc.cronID)

		local, _ := gobase.GetLocalIP()
		_ = metastore.GetInstance().DeleteScheduler(&metastore.ScheduleState{
			JobID: job.Id,
			Host:  local,
		})
	}

	return nil
}

func (s *scheduler) UpdateTaskInstance(ins *sodor.TaskInstance) error {
	ctx, ok := s.jobs.Load(ins.JobId)
	if !ok {
		log.Warnf("cannot found jobCtx. maybe delete already. jobid=%d taskid=%d job_instance=%d task_instance=%d",
			ins.JobId, ins.TaskId, ins.JobInstanceId, ins.Id)
		return errors.New("task is not in a routine job")
	}

	jc := ctx.(*jobContext)
	next, err := jc.UpdateTaskInstance(ins)
	if err != nil {
		return err
	}
	if next > 0 {
		jc.runTask(next)
	}

	return nil
}

func NewParser() cron.Parser {
	return cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
}

type cronLog struct{}

func (l *cronLog) Printf(msg string, v ...interface{}) {
	log.Infof(msg, v...)
}
