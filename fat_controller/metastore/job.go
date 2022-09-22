package metastore

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"gorm.io/gorm"
)

func (ms *metaStore) JobExist(job *sodor.Job) (bool, error) {
	type apiJob struct {
		gorm.Model
	}

	var j Job
	if job.Id > 0 {
		j.ID = uint(job.Id)
	} else if job.Name != "" {
		j.Name = job.Name
	}

	var jobs []apiJob
	rs := ms.db.Where(&j).Find(&jobs)
	if rs.Error != nil {
		return false, rs.Error
	}

	if rs.RowsAffected > 0 {
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

		job.Id = int64(mJob.ID)

		mTasks := make([]Task, 0)
		for _, t := range job.GetTasks() {
			var task Task
			if err := toTask(t, job.Id, &task); err != nil {
				return err
			}
			mTasks = append(mTasks, task)
		}
		if rst := tx.Create(&mTasks); rst.Error != nil {
			return rst.Error
		}

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

			mRels = append(mRels, rel)
		}

		if rst := tx.Create(&mRels); rst.Error != nil {
			return rst.Error
		}

		return nil
	})

	return err
}

func (ms *metaStore) UpdateJob(job *sodor.Job) error {
	gobase.True(job.Id > 0)

	var old sodor.Job
	old.Id = job.Id
	err := ms.SelectJob(&old)
	if err != nil {
		return err
	}

	err = ms.db.Transaction(func(tx *gorm.DB) error {
		var mJob Job
		if err = toJob(job, &mJob); err != nil {
			return err
		}
		if rst := tx.Save(&mJob); rst.Error != nil {
			return rst.Error
		}

		mTasks := make([]Task, 0)
		for _, t := range job.GetTasks() {
			var task Task
			if err = toTask(t, job.Id, &task); err != nil {
				return err
			}
			mTasks = append(mTasks, task)
		}
		if rst := tx.Save(&mTasks); rst.Error != nil {
			return rst.Error
		}

		// Delete the obsolete task the is valid in history
		tasksToDel := make([]int64, 0)
		for _, tOld := range old.GetTasks() {
			found := false
			for _, tNew := range mTasks {
				if tOld.Id == int64(tNew.ID) {
					found = true
					break
				}
			}
			if !found {
				tasksToDel = append(tasksToDel, tOld.Id)
			}
		}
		if len(tasksToDel) > 0 {
			if rs := tx.Delete(&Task{}, tasksToDel); rs.Error != nil {
				return rs.Error
			}
		}

		//tx.Delete(&tasksToDel)

		// Because the relational does not have other associated attributes, it can be deleted directly
		tx.Where("job_id = ?", job.Id).Delete(&TaskRelation{})

		mRels := make([]TaskRelation, 0)
		for _, r := range job.GetRelations() {
			var rel TaskRelation
			rel.JobID = int64(mJob.ID)
			rel.FromTaskID = int64(findTaskID(mTasks, r.FromTask))
			rel.ToTaskID = int64(findTaskID(mTasks, r.ToTask))

			mRels = append(mRels, rel)
		}

		if rst := tx.Create(&mRels); rst.Error != nil {
			return rst.Error
		}

		return nil
	})
	return err
}

func (ms *metaStore) DeleteJob(jID *sodor.Job) error {
	gobase.True(jID.Id > 0)

	return nil
}

// SelectJob returns Job & Tasks & Relations by job.ID
func (ms *metaStore) SelectJob(jID *sodor.Job) error {
	gobase.True(jID.Id > 0)

	var job Job
	rs := ms.db.Find(&job, jID.Id)
	if rs.Error != nil {
		return rs.Error
	}

	if rs.RowsAffected == 0 {
		return ErrNotFound
	}

	if err := fromJob(&job, jID); err != nil {
		return err
	}

	var tasks []Task
	if rs = ms.db.Where(&Task{JobID: int64(job.ID)}).Find(&tasks); rs.Error != nil {
		return rs.Error
	}

	jID.Tasks = make([]*sodor.Task, len(tasks))
	for i, t := range tasks {
		var task sodor.Task
		if err := fromTask(&t, &task); err != nil {
			return err
		}
		jID.Tasks[i] = &task
	}

	var rels []TaskRelation
	rs = ms.db.Where(&TaskRelation{JobID: int64(job.ID)}).Find(&rels)
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

func (ms *metaStore) SelectJobInstance(jID *sodor.Job) error {
	gobase.True(jID.Id > 0)

	return nil
}
