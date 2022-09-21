package metastore

import (
	"github.com/BabySid/proto/sodor"
	"gorm.io/gorm"
)

func (ms *metaStore) InsertJob(job *sodor.Job) error {
	mJob := Job{
		Name:         job.Name,
		AlertRule:    "",
		AlertGroupID: 0,
	}

	mTasks := make([]Task, 0)
	for _, t := range job.GetTasks() {
		var task Task
		task.Name = t.Name
		if t.SchedulerMode == sodor.SchedulerMode_SM_None {
			task.SchedulerMode = int(sodor.SchedulerMode_SM_None.Number())
		} else if t.SchedulerMode == sodor.SchedulerMode_SM_Crontab {
			task.SchedulerMode = int(sodor.SchedulerMode_SM_Crontab.Number())
			task.RoutineSpec = t.RoutineSpec.CtSpec
		}

		task.Script = t.Script
		task.RunningHosts = t.RunningHosts
		task.RunTimeout = int(t.RunningTimeout)

		mTasks = append(mTasks, task)
	}

	return ms.db.Transaction(func(tx *gorm.DB) error {
		rst := tx.Create(&mJob)
		if rst.Error != nil {
			return rst.Error
		}

		rst = tx.Create(&mTasks)
		if rst.Error != nil {
			return rst.Error
		}

		job.Id = int64(mJob.ID)
		taskID := make(map[string]int64)
		for i, t := range mTasks {
			job.Tasks[i].JobId = job.Id
			job.Tasks[i].Id = int64(t.ID)

			taskID[job.Tasks[i].Name] = job.Tasks[i].Id
		}

		mRels := make([]TaskRelation, 0)
		for _, r := range job.GetRelations() {
			var rel TaskRelation
			rel.JobID = int64(mJob.ID)
			rel.FromTaskID = taskID[r.FromTask]
			rel.ToTaskID = taskID[r.ToTask]
			rel.ConditionType = int(r.ConditionType)

			mRels = append(mRels, rel)
		}

		rst = tx.Create(&mRels)
		return nil
	})
}

func (ms *metaStore) UpdateJob() error {
	return nil
}

func (ms *metaStore) DeleteJob() error {
	return nil
}

func (ms *metaStore) SelectJob() error {
	return nil
}
