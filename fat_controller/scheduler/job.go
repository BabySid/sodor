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

	lastJobInsID int32
	// jobInsID =>
	instances map[int32]*instance

	// alertPluginInstanceID => alert
	alerts map[int32]alert.Alert
}

type instance struct {
	jobInstance *sodor.JobInstance
	// taskID=>taskInsID =>
	taskInstances map[int32]map[int32]*sodor.TaskInstance
}

func newJobContext() *jobContext {
	return &jobContext{
		job:          nil,
		jobDag:       nil,
		cronID:       0,
		lastJobInsID: 0,
		instances:    make(map[int32]*instance),
		alerts:       nil,
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
			ding := alert.NewDingDing(plugin.Dingding.Webhook, plugin.Dingding.Sign, plugin.Dingding.AtMobiles)
			jc.alerts[int32(id)] = ding
		}
	}

	return nil
}

func (jc *jobContext) Run() {
	jc.lock.Lock()
	defer jc.lock.Unlock()

	if jc.lastJobInsID > 0 {
		if ins, ok := jc.instances[jc.lastJobInsID]; ok {
			if ins.jobInstance != nil && ins.jobInstance.StopTs == 0 {
				// todo schedule-strategy: skip_run_when_last_is_running
				logJob(jc.job).Infof("run job delayed because of last instance(%d) is not done.", jc.lastJobInsID)
				return
			}
		}
	}

	curInstance := &sodor.JobInstance{
		JobId:      jc.job.Id,
		ScheduleTs: int32(time.Now().Unix()),
		StartTs:    int32(time.Now().Unix()),
	}

	// to support broadcast, we need generate a new task instance according to running_host
	// and fat_ctrl update the task instance with the reply from thomas's host and instance_id
	taskInstances := make([]*sodor.TaskInstance, 0)
	for _, t := range jc.job.Tasks {
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
		for _, h := range t.RunningHosts {
			var ti sodor.TaskInstance
			ti.Host = h.Node
			ti.TaskId = taskIns.TaskId
			ti.JobId = taskIns.JobId
			ti.StartTs = taskIns.StartTs
			ti.ParsedContent = taskIns.ParsedContent
			taskInstances = append(taskInstances, &ti)
		}
	}

	if err := metastore.GetInstance().InsertJobTaskInstance(curInstance, taskInstances); err != nil {
		logJob(jc.job).Warnf("run job failed. InsertJobTaskInstance return err=%s", err)
		alert.GetInstance().GiveAlert(fmt.Sprintf("run job failed. InsertJobTaskInstance return err=%s", err))
		return
	}

	jc.lastJobInsID = curInstance.Id

	taskInsMap := make(map[int32]map[int32]*sodor.TaskInstance)
	for _, t := range taskInstances {
		if v, ok := taskInsMap[t.TaskId]; ok {
			v[t.Id] = t
		} else {
			insMap := make(map[int32]*sodor.TaskInstance)
			insMap[t.Id] = t
			taskInsMap[t.TaskId] = insMap
		}
	}
	jc.instances[curInstance.Id] = &instance{
		jobInstance:   curInstance,
		taskInstances: taskInsMap,
	}

	logJob(jc.job).Infof("build task instances. jobInsId=%d sizeOfTaskIns=%d", curInstance.Id, len(taskInstances))
	for tid, tIns := range jc.instances[curInstance.Id].taskInstances {
		for _, ins := range tIns {
			logJob(jc.job).Infof("taskID:%d host:%s taskInsID:%d", tid, ins.Host, ins.Id)
		}
	}

	go func() {
		firstTask := int32(jc.jobDag.topoNodes[0].ID())
		logJob(jc.job).Infof("begin to run job. jobInsId=%d firstTaskID=%d", curInstance.Id, firstTask)
		taskIns := jc.findTaskInstance(curInstance.Id, firstTask)
		task := jc.findTask(firstTask)
		jc.runTask(task, taskIns)
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
		taskInsMap := make(map[int32]map[int32]*sodor.TaskInstance)
		for _, t := range taskInstances.TaskInstances {
			if v, ok := taskInsMap[t.TaskId]; ok {
				v[t.Id] = t
			} else {
				insMap := make(map[int32]*sodor.TaskInstance)
				insMap[t.Id] = t
				taskInsMap[t.TaskId] = insMap
			}
		}

		instances = &instance{
			jobInstance:   &curInstance,
			taskInstances: taskInsMap,
		}
		jc.instances[ins.JobInstanceId] = instances
	}

	taskDone := true
	// update job_instance & task_instance
	if v, ok := instances.taskInstances[ins.TaskId]; ok {
		v[ins.Id] = ins
		for _, is := range v {
			if is.StopTs == 0 || is.ExitCode != 0 {
				taskDone = false
			}
		}
	}

	nextTask := 0

	if ins.ExitCode != 0 {
		jc.buildJobInstance(ins, instances.jobInstance)
	} else {
		if taskDone {
			for i, node := range jc.jobDag.topoNodes {
				if node.ID() == int64(ins.TaskId) {
					if i == len(jc.jobDag.topoNodes)-1 {
						jc.buildJobInstance(ins, instances.jobInstance)
					} else {
						nextTask = int(jc.jobDag.topoNodes[i+1].ID())
					}
					break
				}
			}
		}
	}

	logJob(jc.job).Infof("UpdateTaskInstance(taskInsId:%d) from host:%s with stopts=%d exit_code:%d nextTask:%d taskDone:%v",
		ins.TaskId, ins.Host, ins.StopTs, ins.ExitCode, nextTask, taskDone)

	var err error
	if nextTask != 0 {
		err = metastore.GetInstance().UpdateJobTaskInstance(nil, ins)
	} else if ins.ExitCode != 0 {
		err = metastore.GetInstance().UpdateJobTaskInstance(instances.jobInstance, ins)
		if instances.jobInstance.ExitCode != 0 {
			msg := fmt.Sprintf("job:%s finished with a error:%s from task:%s",
				jc.job.Name, instances.jobInstance.ExitMsg, jc.findTask(ins.TaskId).Name)
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

func (jc *jobContext) findTaskInstance(jobIns int32, taskId int32) *sodor.TaskInstances {
	jc.lock.Lock()
	defer jc.lock.Unlock()

	var taskIns sodor.TaskInstances
	taskIns.TaskInstances = make([]*sodor.TaskInstance, 0)

	for _, ins := range jc.instances[jobIns].taskInstances[taskId] {
		taskIns.TaskInstances = append(taskIns.TaskInstances, ins)
	}

	gobase.True(len(taskIns.TaskInstances) > 0)
	return &taskIns
}

func (jc *jobContext) runTask(task *sodor.Task, taskInstances *sodor.TaskInstances) {
	var err error
	var th *metastore.Thomas
	var ins *sodor.TaskInstance
	defer func() {
		if err != nil {
			jc.terminalJob(task, ins, err)
		}
	}()

	for _, instance := range taskInstances.TaskInstances {
		ins = instance
		th, err = metastore.GetInstance().SelectValidThomas(ins.Host)
		if err != nil {
			return
		}

		if th.ID == 0 {
			err = fmt.Errorf("no thomas found for host(%s)", ins.Host)
			return
		}

		err = jc.sendTaskToThomas(th, task, ins)
		logJob(jc.job).Infof("run task(%d:%s) at host:%s. taskInsID:%d err=%v", task.Id, task.Name, ins.Host, ins.Id, err)
	}
}

func (jc *jobContext) terminalJob(task *sodor.Task, ins *sodor.TaskInstance, cause error) {
	insMap := jc.instances[ins.JobInstanceId].taskInstances[ins.TaskId]
	gobase.True(insMap != nil)
	taskIns := insMap[ins.Id]
	gobase.True(taskIns != nil)

	taskIns.Id = ins.Id
	taskIns.JobId = task.JobId
	taskIns.TaskId = task.Id
	taskIns.JobInstanceId = ins.JobInstanceId
	taskIns.StopTs = int32(time.Now().Unix())
	taskIns.Host = ins.Host
	taskIns.ExitCode = -1
	taskIns.ExitMsg = cause.Error()

	logJob(jc.job).Infof("terminal task(%d:%s) with error:%s. taskInsID:%d", task.Id, task.Name, taskIns.ExitMsg, ins.Id)

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
