package scheduler

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/thomas"
	"sync"
	"time"
)

type jobContext struct {
	lock   sync.Mutex
	job    *sodor.Job
	jobDag *dag
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
		job:       nil,
		jobDag:    nil,
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
		StartTs:    int32(time.Now().Unix()),
	}

	taskInstances := make([]*sodor.TaskInstance, len(jc.job.Tasks))
	for i, t := range jc.job.Tasks {
		var taskIns sodor.TaskInstance
		taskIns.TaskId = t.Id
		taskIns.JobId = t.JobId
		taskIns.StartTs = int32(time.Now().Unix())
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

	go func() {
		taskIns, task := jc.getTaskInstance(curInstance.Id, int32(jc.jobDag.topoNodes[0].ID()))
		jc.runTask(curInstance.Id, taskIns, task)
		logJob(jc.job).Infof("run job at %s, job_instance_id=%d", gobase.FormatTimeStamp(int64(curInstance.ScheduleTs)), curInstance.Id)
	}()

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
			log.Warnf("loadInstanceFromMetaStore return err=%v. taskInstance.Id=%d jobId=%d taskId=%d",
				err, ins.Id, ins.JobId, ins.TaskId)
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
				break
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

func (jc *jobContext) getTaskInstance(jobIns int32, taskId int32) (int32, *sodor.Task) {
	jc.lock.Lock()
	defer jc.lock.Unlock()

	var taskIns int32
	var task *sodor.Task

	for id, t := range jc.instances[jobIns].taskInstances {
		if t.TaskId == taskId {
			taskIns = id
		}
	}
	for _, t := range jc.job.Tasks {
		if t.Id == taskId {
			task = t
		}
	}

	gobase.True(taskIns > 0 && task != nil)
	return taskIns, task
}

func (jc *jobContext) runTask(jobIns int32, taskIns int32, task *sodor.Task) {
	var err error
	var th *metastore.Thomas
	defer func() {
		if err != nil {
			jc.terminalJob(jobIns, taskIns, task, err)
		}
		log.Infof("run task task_id=%d task_name=%s err=%v", task.Id, task.Name, err)
	}()

	th, err = metastore.GetInstance().SelectValidThomas(task.RunningHosts[0].Node)
	if err != nil {
		return
	}

	if th.ID == 0 {
		err = fmt.Errorf("no thomas found for host(%s)", task.RunningHosts[0].Node)
		return
	}

	err = jc.sendTaskToThomas(th, jobIns, taskIns, task)
}

func (jc *jobContext) terminalJob(jobInsID int32, taskInsID int32, task *sodor.Task, cause error) {
	taskIns := jc.instances[jobInsID].taskInstances[taskInsID]
	taskIns.Id = taskInsID
	taskIns.JobId = task.JobId
	taskIns.TaskId = task.Id
	taskIns.JobInstanceId = jobInsID
	taskIns.StopTs = int32(time.Now().Unix())
	taskIns.Host = task.RunningHosts[0].Node
	taskIns.ExitCode = -1
	taskIns.ExitMsg = cause.Error()

	_, err := jc.UpdateTaskInstance(taskIns)
	if err != nil {
		log.Warnf("UpdateTaskInstance failed. err=%s", err)
	}
}

func (jc *jobContext) loadInstanceFromMetaStore(job *sodor.JobInstance, task *sodor.TaskInstances) error {
	err := metastore.GetInstance().SelectJobTaskInstance(job, task)
	return err
}

func (jc *jobContext) sendTaskToThomas(th *metastore.Thomas, jobIns int32, taskIns int32, task *sodor.Task) error {
	t := thomas.Thomas{
		Host: th.Host,
		Port: th.Port,
	}
	return t.RunTask(jobIns, taskIns, task)
}

func logJob(jc *sodor.Job) *log.Entry {
	return log.WithFields(log.Fields{"id": jc.Id, "name": jc.Name})
}
