package metastore

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"gorm.io/gorm"
)

func (ms *metaStore) InsertJobTaskInstance(job *sodor.JobInstance, tasks []*sodor.TaskInstance) error {
	gobase.True(job.Id == 0)

	err := ms.db.Transaction(func(tx *gorm.DB) error {
		var ins JobInstance
		if err := toJobInstance(job, &ins); err != nil {
			return err
		}

		if rs := ms.db.Create(&ins); rs.Error != nil {
			return rs.Error
		}

		job.Id = int32(ins.ID)

		for _, t := range tasks {
			var ins TaskInstance
			if err := toTaskInstance(t, &ins); err != nil {
				return err
			}

			if rs := ms.db.Create(&ins); rs.Error != nil {
				return rs.Error
			}

			t.Id = int32(ins.ID)
			t.JobInstanceId = job.Id
		}

		return nil
	})

	return err
}

func (ms *metaStore) UpdateJobTaskInstance(job *sodor.JobInstance, task *sodor.TaskInstance) error {
	err := ms.db.Transaction(func(tx *gorm.DB) error {
		if job != nil {
			gobase.True(job.Id > 0)
			var ins JobInstance
			if err := toJobInstance(job, &ins); err != nil {
				return err
			}

			if rst := ms.db.Save([]*JobInstance{&ins}); rst.Error != nil {
				return rst.Error
			}
		}

		var ins TaskInstance
		if err := toTaskInstance(task, &ins); err != nil {
			return err
		}

		if rst := ms.db.Save([]*TaskInstance{&ins}); rst.Error != nil {
			return rst.Error
		}
		return nil
	})

	return err
}

func (ms *metaStore) SelectJobTaskInstance(job *sodor.JobInstance, tasks *sodor.TaskInstances) error {
	gobase.True(job.Id > 0)

	var jobIns JobInstance
	rs := ms.db.Limit(1).Find(&jobIns, job.Id)
	if rs.Error != nil {
		return rs.Error
	}

	if rs.RowsAffected == 0 {
		return ErrNotFound
	}

	if err := fromJobInstance(&jobIns, job); err != nil {
		return err
	}

	var taskIns []*TaskInstance
	if rs = ms.db.Where(&TaskInstance{JobInstanceID: job.Id}).Find(&taskIns); rs.Error != nil {
		return rs.Error
	}

	tasks.TaskInstances = make([]*sodor.TaskInstance, len(taskIns))
	for i, t := range taskIns {
		var task sodor.TaskInstance
		if err := fromTaskInstance(t, &task); err != nil {
			return err
		}
		tasks.TaskInstances[i] = &task
	}

	return nil
}
