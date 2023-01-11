package metastore

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"gorm.io/gorm"
	"sodor/base"
)

func (ms *metaStore) JobExist(job *sodor.Job) (bool, error) {
	type apiJob struct {
		gobase.TableModel
	}

	var j Job
	if job.Id > 0 {
		j.ID = uint(job.Id)
	} else if job.Name != "" {
		j.Name = job.Name
	}

	var jobs apiJob
	rs := ms.db.Model(&Job{}).Where(&j).Limit(1).Find(&jobs)
	if rs.Error != nil {
		return false, rs.Error
	}

	if jobs.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (ms *metaStore) InsertJob(job *sodor.Job) error {
	gobase.True(job.Id == 0)
	err := ms.db.Transaction(func(tx *gorm.DB) error {
		var mJob Job
		if err := toJob(job, &mJob); err != nil {
			return err
		}

		if rst := tx.Create(&mJob); rst.Error != nil {
			return rst.Error
		}

		job.Id = int32(mJob.ID)

		mTasks := make([]Task, len(job.GetTasks()))
		for i, t := range job.GetTasks() {
			var task Task
			if err := toTask(t, job.Id, &task); err != nil {
				return err
			}
			mTasks[i] = task
		}
		if rst := tx.Create(&mTasks); rst.Error != nil {
			return rst.Error
		}

		taskID := make(map[string]int32)
		for i, t := range mTasks {
			job.Tasks[i].JobId = job.Id
			job.Tasks[i].Id = int32(t.ID)

			taskID[job.Tasks[i].Name] = job.Tasks[i].Id
		}

		mRels := make([]TaskRelation, 0)
		for _, r := range job.GetRelations() {
			var rel TaskRelation
			rel.JobID = int32(mJob.ID)
			rel.FromTaskID = taskID[r.FromTask]
			rel.ToTaskID = taskID[r.ToTask]

			mRels = append(mRels, rel)
		}

		if len(mRels) > 0 {
			if rst := tx.Create(&mRels); rst.Error != nil {
				return rst.Error
			}
		}

		if job.ScheduleMode == sodor.ScheduleMode_ScheduleMode_Crontab {
			stat := &ScheduleState{
				JobID: job.Id,
				Host:  base.LocalHost,
			}

			if rst := tx.Create(&stat); rst.Error != nil {
				return rst.Error
			}
		}

		return nil
	})

	return err
}

func (ms *metaStore) UpdateJob(job *sodor.Job) error {
	gobase.True(job.Id > 0)

	err := ms.db.Transaction(func(tx *gorm.DB) error {
		// 1. update job meta
		var mJob Job
		if err := toJob(job, &mJob); err != nil {
			return err
		}

		if rst := tx.Model(&mJob).Select(mJob.UpdateFields()).Updates(mJob); rst.Error != nil {
			return rst.Error
		}

		// 2. update task
		newTasks := make([]*Task, 0)
		updatedTasks := make([]*Task, 0)
		for _, t := range job.GetTasks() {
			var task Task
			if err := toTask(t, job.Id, &task); err != nil {
				return err
			}
			if t.Id == 0 {
				newTasks = append(newTasks, &task)
			} else {
				updatedTasks = append(updatedTasks, &task)
			}
		}

		tasksToRetain := make([]uint, 0)
		if len(newTasks) > 0 {
			if rst := tx.Create(&newTasks); rst.Error != nil {
				return rst.Error
			}

			for _, t := range newTasks {
				tasksToRetain = append(tasksToRetain, t.ID)
			}
		}

		for _, task := range updatedTasks {
			if rst := tx.Model(task).Select(task.UpdateFields()).Updates(task); rst.Error != nil {
				return rst.Error
			}

			tasksToRetain = append(tasksToRetain, task.ID)
		}

		// Delete the obsolete task that is valid in history
		if rs := tx.Where("job_id = ? and id not in ?", job.Id, tasksToRetain).Delete(&Task{}); rs.Error != nil {
			return rs.Error
		}

		// 3. update relation
		// Because the relational does not have other associated attributes, it can be deleted directly
		if rs := tx.Where("job_id = ?", job.Id).Delete(&TaskRelation{}); rs.Error != nil {
			return rs.Error
		}

		mRels := make([]TaskRelation, 0)
		allTasks := make([]*Task, 0)
		allTasks = append(append(allTasks, newTasks...), updatedTasks...)
		for _, r := range job.GetRelations() {
			var rel TaskRelation
			rel.JobID = int32(mJob.ID)
			rel.FromTaskID = int32(findTaskID(allTasks, r.FromTask))
			rel.ToTaskID = int32(findTaskID(allTasks, r.ToTask))

			mRels = append(mRels, rel)
		}

		if len(mRels) > 0 {
			if rst := tx.Create(&mRels); rst.Error != nil {
				return rst.Error
			}
		}

		// 4. update scheduler state
		if rst := tx.Where(&ScheduleState{JobID: job.Id, Host: base.LocalHost}).Delete(ScheduleState{}); rst.Error != nil {
			return rst.Error
		}

		if job.ScheduleMode == sodor.ScheduleMode_ScheduleMode_Crontab {
			stat := &ScheduleState{
				JobID: job.Id,
				Host:  base.LocalHost,
			}

			if rst := tx.Create(&stat); rst.Error != nil {
				return rst.Error
			}
		}

		return nil
	})
	return err
}

func (ms *metaStore) DeleteJob(jID *sodor.Job) error {
	gobase.True(jID.Id > 0)

	err := ms.db.Transaction(func(tx *gorm.DB) error {
		var job Job
		job.ID = uint(jID.Id)

		if rs := tx.Where(job).Delete(&Job{}); rs.Error != nil {
			return rs.Error
		}

		if rs := tx.Where(Task{JobID: int32(job.ID)}).Delete(&Task{}); rs.Error != nil {
			return rs.Error
		}

		if rs := tx.Where(TaskRelation{JobID: int32(job.ID)}).Delete(&TaskRelation{}); rs.Error != nil {
			return rs.Error
		}

		if rs := tx.Where(&ScheduleState{JobID: jID.Id, Host: base.LocalHost}).Delete(ScheduleState{}); rs.Error != nil {
			return rs.Error
		}

		if rs := tx.Where(JobInstance{JobID: int32(job.ID)}).Delete(&JobInstance{}); rs.Error != nil {
			return rs.Error
		}

		if rs := tx.Where(TaskInstance{JobID: int32(job.ID)}).Delete(&TaskInstance{}); rs.Error != nil {
			return rs.Error
		}

		return nil
	})

	return err
}

// SelectJob returns Job & Tasks & Relations by job.ID
func (ms *metaStore) SelectJob(jID *sodor.Job) error {
	gobase.True(jID.Id > 0)

	var job Job
	rs := ms.db.Limit(1).Find(&job, jID.Id)
	if rs.Error != nil {
		return rs.Error
	}

	if rs.RowsAffected == 0 {
		return ErrNotFound
	}

	if err := fromJob(&job, jID); err != nil {
		return err
	}

	var tasks []*Task
	if rs = ms.db.Where(&Task{JobID: int32(job.ID)}).Find(&tasks); rs.Error != nil {
		return rs.Error
	}

	jID.Tasks = make([]*sodor.Task, len(tasks))
	for i, t := range tasks {
		var task sodor.Task
		if err := fromTask(t, &task); err != nil {
			return err
		}
		jID.Tasks[i] = &task
	}

	var rels []TaskRelation
	rs = ms.db.Where(&TaskRelation{JobID: int32(job.ID)}).Find(&rels)
	if rs.Error != nil {
		return rs.Error
	}

	jID.Relations = make([]*sodor.TaskRelation, len(rels))
	for i, r := range rels {
		var rel sodor.TaskRelation
		rel.FromTask = findTaskName(tasks, uint(r.FromTaskID))
		rel.ToTask = findTaskName(tasks, uint(r.ToTaskID))

		jID.Relations[i] = &rel
	}

	return nil
}

func (ms *metaStore) ListJobs() (*sodor.Jobs, error) {
	var jobs []Job
	if rst := ms.db.Find(&jobs); rst.Error != nil {
		return nil, rst.Error
	}

	var sJobs sodor.Jobs
	sJobs.Jobs = make([]*sodor.Job, len(jobs))
	for i, job := range jobs {
		var sJob sodor.Job
		err := fromJob(&job, &sJob)
		if err != nil {
			return nil, err
		}

		sJobs.Jobs[i] = &sJob
	}

	return &sJobs, nil
}
