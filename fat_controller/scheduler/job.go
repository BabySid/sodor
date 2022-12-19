package scheduler

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sodor/fat_controller/alert"
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
	// jobInsID =>
	instances map[int32]*instance

	// alertPluginInstanceID => alert
	alerts map[int32]alert.Alert
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
		alerts:    nil,
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
	err = jc.setAlerts()
	if err != nil {
		return err
	}
	return nil
}

func (jc *jobContext) setAlerts() error {
	if jc.job.AlertGroupId > 0 {
		ag := sodor.AlertGroup{
			Id: jc.job.AlertGroupId,
		}
		plugins := sodor.AlertPluginInstances{}

		err := metastore.GetInstance().ShowAlertGroup(&ag, &plugins)
		if err != nil {
			return err
		}

		for id, plugin := range plugins.AlertPluginInstances {
			param := plugin.Plugin.(*sodor.AlertPluginInstance_Dingding)
			ding := alert.NewDingDing(param.Dingding.Webhook, param.Dingding.Sign, param.Dingding.AtMobiles)
			jc.alerts[int32(id)] = ding
		}
	}

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

	// todo to support broadcast, we need generate a new task instance according to running_host
	// and fat_ctrl update the task instance with the reply from thomas's host and instance_id
	taskInstances := make([]*sodor.TaskInstance, len(jc.job.Tasks))
	for i, t := range jc.job.Tasks {
		var taskIns sodor.TaskInstance
		taskIns.TaskId = t.Id
		taskIns.JobId = t.JobId
		taskIns.StartTs = int32(time.Now().Unix())
		// todo parse the content according task_type
		if err := parseTaskContent(t, &taskIns); err != nil {
			logJob(jc.job).Warnf("parseTaskContent for job failed. err=%s", err)
			alert.GetInstance().GiveAlert(fmt.Sprintf("parseTaskContent for job failed. err=%s", err))
			return
		}
		taskInstances[i] = &taskIns
	}
	if err := metastore.GetInstance().InsertJobTaskInstance(curInstance, taskInstances); err != nil {
		logJob(jc.job).Warnf("run job failed. InsertJobTaskInstance return err=%s", err)
		alert.GetInstance().GiveAlert(fmt.Sprintf("run job failed. InsertJobTaskInstance return err=%s", err))
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
		taskIns := jc.findTaskInstance(curInstance.Id, int32(jc.jobDag.topoNodes[0].ID()))
		task := jc.findTask(int32(jc.jobDag.topoNodes[0].ID()))
		jc.runTask(task, taskIns)
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
	if _, ok := instances.taskInstances[ins.Id]; ok {
		instances.taskInstances[ins.Id] = ins
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
	if nextTask != 0 {
		err = metastore.GetInstance().UpdateJobTaskInstance(nil, ins)
	} else {
		err = metastore.GetInstance().UpdateJobTaskInstance(instances.curInstance, ins)
		if instances.curInstance.ExitCode != 0 {
			msg := fmt.Sprintf("job:%s finished with a error:%s from task:%s",
				jc.job.Name, instances.curInstance.ExitMsg, jc.findTask(ins.TaskId).Name)
			jc.giveAlert(msg)
		}
	}

	return int32(nextTask), err
}

func (jc *jobContext) buildJobInstance(ins *sodor.TaskInstance, jobIns *sodor.JobInstance) {
	jobIns.StopTs = ins.StopTs
	jobIns.ExitCode = ins.ExitCode
	jobIns.ExitMsg = ins.ExitMsg
}

func (jc *jobContext) findTask(taskId int32) *sodor.Task {
	jc.lock.Lock()
	defer jc.lock.Unlock()

	for _, t := range jc.job.Tasks {
		if t.Id == taskId {
			return t
		}
	}

	gobase.AssertHere()
	return nil
}

func (jc *jobContext) findTaskInstance(jobIns int32, taskId int32) *sodor.TaskInstance {
	jc.lock.Lock()
	defer jc.lock.Unlock()

	var taskIns *sodor.TaskInstance

	for _, t := range jc.instances[jobIns].taskInstances {
		if t.TaskId == taskId {
			taskIns = t
		}
	}

	gobase.True(taskIns != nil)
	return taskIns
}

func (jc *jobContext) runTask(task *sodor.Task, ins *sodor.TaskInstance) {
	var err error
	var th *metastore.Thomas
	defer func() {
		if err != nil {
			jc.terminalJob(task, ins, err)
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

	err = jc.sendTaskToThomas(th, task, ins)
}

func (jc *jobContext) terminalJob(task *sodor.Task, ins *sodor.TaskInstance, cause error) {
	taskIns := jc.instances[ins.JobInstanceId].taskInstances[ins.Id]
	taskIns.Id = ins.Id
	taskIns.JobId = task.JobId
	taskIns.TaskId = task.Id
	taskIns.JobInstanceId = ins.JobInstanceId
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
	return metastore.GetInstance().SelectInstanceByJobInsID(job, task)
}

func (jc *jobContext) sendTaskToThomas(th *metastore.Thomas, task *sodor.Task, ins *sodor.TaskInstance) error {
	t := thomas.Thomas{
		Host: th.Host,
		Port: th.Port,
	}
	return t.RunTask(task, ins)
}

func (jc *jobContext) giveAlert(msg string) {
	for id, v := range jc.alerts {
		err := v.GiveAlarm(msg)
		status := "OK"
		if err != nil {
			status = err.Error()
		}
		his := sodor.AlertPluginInstanceHistory{
			InstanceId: id,
			GroupId:    jc.job.AlertGroupId,
			AlertMsg:   msg,
			StatusMsg:  status,
		}
		err = metastore.GetInstance().InsertAlertPluginInstanceHistory(&his)
		if err != nil {
			logJob(jc.job).Warnf("giveAlert failed. pluginInstanceID=%d err=%s", id, err)
		}
	}
}

func logJob(jc *sodor.Job) *log.Entry {
	return log.WithFields(log.Fields{"id": jc.Id, "name": jc.Name})
}
