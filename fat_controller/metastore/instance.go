package metastore

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"gorm.io/gorm"
	"sodor/fat_controller/config"
)

func (ms *metaStore) InsertJobTaskInstance(job *sodor.JobInstance, tasks []*sodor.TaskInstance) error {
	gobase.True(job.Id == 0)

	err := ms.db.Transaction(func(tx *gorm.DB) error {
		var ins JobInstance
		if err := toJobInstance(job, &ins); err != nil {
			return err
		}

		if rs := tx.Create(&ins); rs.Error != nil {
			return rs.Error
		}

		job.Id = int32(ins.ID)

		for _, t := range tasks {
			t.JobInstanceId = job.Id

			var ins TaskInstance
			if err := toTaskInstance(t, &ins); err != nil {
				return err
			}

			if rs := tx.Create(&ins); rs.Error != nil {
				return rs.Error
			}

			t.Id = int32(ins.ID)
		}

		// delete long-age records
		if config.GetInstance().MaxJobInstance <= 0 {
			return nil
		}

		var jobIds []uint
		rs := tx.Model(&JobInstance{}).Where(JobInstance{JobID: job.JobId}).
			Order("id desc").Offset(int(config.GetInstance().MaxJobInstance)).Limit(1024).
			Pluck("id", &jobIds)
		if rs.Error != nil {
			return rs.Error
		}

		if len(jobIds) == 0 {
			return nil
		}

		if rs = tx.Where("id in (?)", jobIds).Delete(&JobInstance{}); rs.Error != nil {
			return rs.Error
		}

		if rs = tx.Where("job_instance_id in (?)", jobIds).Delete(&TaskInstance{}); rs.Error != nil {
			return rs.Error
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

			if rst := tx.Model(&ins).Select(ins.UpdateFields()).Updates(ins); rst.Error != nil {
				return rst.Error
			}
		}

		var ins TaskInstance
		if err := toTaskInstance(task, &ins); err != nil {
			return err
		}

		if rst := tx.Model(&ins).Select(ins.UpdateFields()).Updates(ins); rst.Error != nil {
			return rst.Error
		}
		return nil
	})

	return err
}

func (ms *metaStore) SelectInstanceByJobInsID(job *sodor.JobInstance, tasks *sodor.TaskInstances) error {
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

func (ms *metaStore) SelectInstanceByJobID(jobID int32) (*sodor.JobTaskInstances, error) {
	var jobIns []*JobInstance
	rs := ms.db.Where(&JobInstance{JobID: jobID}).Find(&jobIns)
	if rs.Error != nil {
		return nil, rs.Error
	}

	var taskIns []*TaskInstance
	rs = ms.db.Where(&TaskInstance{JobID: jobID}).Find(&taskIns)
	if rs.Error != nil {
		return nil, rs.Error
	}

	var jtIns sodor.JobTaskInstances
	jtIns.JobInstances.JobInstances = make([]*sodor.JobInstance, 0)
	jtIns.TaskInstances.TaskInstances = make([]*sodor.TaskInstance, 0)

	for _, ins := range jobIns {
		var target sodor.JobInstance
		err := fromJobInstance(ins, &target)
		if err != nil {
			return nil, err
		}

		jtIns.JobInstances.JobInstances = append(jtIns.JobInstances.JobInstances, &target)
	}

	for _, ins := range taskIns {
		var target sodor.TaskInstance
		err := fromTaskInstance(ins, &target)
		if err != nil {
			return nil, err
		}

		jtIns.TaskInstances.TaskInstances = append(jtIns.TaskInstances.TaskInstances, &target)
	}

	return &jtIns, nil
}
