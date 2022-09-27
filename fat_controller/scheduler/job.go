package scheduler

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sodor/fat_controller/metastore"
	"sync"
	"time"
)

type jobContext struct {
	lock   sync.Mutex
	job    *sodor.Job
	jobDag *dag
	//dirty  bool // True indicates that job needs to be synchronized from the metastore layer
	cronID cron.EntryID

	lastInsID int32
	// jobInsID = >
	instances map[int32]*instance
}

type instance struct {
	curInstance *sodor.JobInstance
	// taskInsID =>
	taskInstances map[int32]*sodor.TaskInstance
}

func newJobContext() *jobContext {
	return &jobContext{
		job:    nil,
		jobDag: nil,
		//dirty:  true,
		cronID:    0,
		lastInsID: 0,
		instances: make(map[int32]*instance),
	}
}

func (jc *jobContext) setJob(j *sodor.Job) error {
	jc.lock.Lock()
	defer jc.lock.Unlock()

	g := newDAG()
	err := g.buildFromJob(j)
	if err != nil {
		return err
	}
	jc.job = j
	jc.jobDag = g
	return nil
}

func (jc *jobContext) Run() {
	jc.lock.Lock()
	defer jc.lock.Unlock()

	if jc.lastInsID > 0 && jc.instances[jc.lastInsID].curInstance != nil && jc.instances[jc.lastInsID].curInstance.StopTs == 0 {
		// todo schedule-strategy: skip_run_when_last_is_running
		logJob(jc.job).Infof("run job delayed because of last instance(%d) is not done.", jc.lastInsID)
		return
	}

	curInstance := &sodor.JobInstance{
		JobId:      jc.job.Id,
		ScheduleTs: int32(time.Now().Unix()),
	}

	taskInstances := make([]*sodor.TaskInstance, len(jc.job.Tasks))
	for i, t := range jc.job.Tasks {
		var taskIns sodor.TaskInstance
		taskIns.TaskId = t.Id
		taskIns.JobId = t.JobId
		taskInstances[i] = &taskIns
	}
	if err := metastore.GetInstance().InsertJobTaskInstance(curInstance, taskInstances); err != nil {
		// todo system warning
		logJob(jc.job).Warnf("run job failed. InsertJobTaskInstance return err=%s", err)
		return
	}

	taskInsMap := make(map[int32]*sodor.TaskInstance)
	for _, t := range taskInstances {
		taskInsMap[t.Id] = t
	}
	jc.instances[curInstance.Id] = &instance{
		curInstance:   curInstance,
		taskInstances: taskInsMap,
	}

	jc.runTask(int32(jc.jobDag.topoNodes[0].ID()))
	logJob(jc.job).Infof("run job at %s, instance_id=%d", gobase.FormatTimeStamp(int64(curInstance.ScheduleTs)), curInstance.Id)
	return
}

func (jc *jobContext) UpdateTaskInstance(ins *sodor.TaskInstance) (int32, error) {
	jc.lock.Lock()
	defer jc.lock.Unlock()

	instances := jc.instances[ins.JobInstanceId]
	if instances == nil {
		var curInstance sodor.JobInstance
		curInstance.Id = ins.JobInstanceId
		var taskInstances sodor.TaskInstances
		err := jc.loadInstanceFromMetaStore(&curInstance, &taskInstances)
		if err != nil {
			return 0, err
		}
		taskInsMap := make(map[int32]*sodor.TaskInstance)
		for _, t := range taskInstances.TaskInstances {
			taskInsMap[t.Id] = t
		}

		instances = &instance{
			curInstance:   &curInstance,
			taskInstances: taskInsMap,
		}
		jc.instances[ins.JobInstanceId] = instances
	}

	// update job_instance & task_instance
	if task, ok := instances.taskInstances[ins.Id]; ok {
		*task = *ins
	}

	nextTask := 0

	if ins.ExitCode != 0 {
		jc.buildJobInstance(ins, instances.curInstance)
	} else {
		for i, node := range jc.jobDag.topoNodes {
			if node.ID() == int64(ins.TaskId) {
				if i == len(jc.jobDag.topoNodes)-1 {
					jc.buildJobInstance(ins, instances.curInstance)
				} else {
					nextTask = int(jc.jobDag.topoNodes[i+1].ID())
				}
			}
		}
	}

	var err error
	if nextTask == 0 {
		err = metastore.GetInstance().UpdateJobTaskInstance(nil, ins)
	} else {
		err = metastore.GetInstance().UpdateJobTaskInstance(instances.curInstance, ins)
	}

	return int32(nextTask), err
}

func (jc *jobContext) buildJobInstance(ins *sodor.TaskInstance, jobIns *sodor.JobInstance) {
	jobIns.StopTs = ins.StopTs
	jobIns.ExitCode = ins.ExitCode
	jobIns.ExitMsg = ins.ExitMsg
}

func (jc *jobContext) runTask(taskId int32) {
	// select nodes from thomas
	var task *sodor.Task
	for _, t := range jc.job.Tasks {
		if t.Id == taskId {
			task = t
		}
	}

	gobase.True(task != nil)
	log.Infof("run task task_id=%d task_name=%s", task.Id, task.Name)
	// send task request to thomas
	// if error found in scheduler, then call UpdateTaskInstance(...)
}

func (jc *jobContext) loadInstanceFromMetaStore(job *sodor.JobInstance, task *sodor.TaskInstances) error {
	err := metastore.GetInstance().SelectJobTaskInstance(job, task)
	return err
}

func logJob(jc *sodor.Job) *log.Entry {
	return log.WithFields(log.Fields{"id": jc.Id, "name": jc.Name})
}
