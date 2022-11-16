package metastore

import (
	"github.com/BabySid/gobase"
	"gorm.io/gorm"
)

func (ms *metaStore) SelectScheduler(host string) ([]ScheduleState, error) {
	var stats []ScheduleState
	if rs := ms.db.Where(&ScheduleState{Host: host}).Find(&stats); rs.Error != nil {
		return nil, rs.Error
	}
	return stats, nil
}

func (ms *metaStore) UpsertScheduler(stat *ScheduleState) error {
	gobase.True(stat.JobID > 0 && stat.Host != "")

	err := ms.db.Transaction(func(tx *gorm.DB) error {
		var s ScheduleState
		if rs := ms.db.Where(stat).Limit(1).Find(&s); rs.Error != nil {
			return rs.Error
		}

		if s.ID > 0 {
			if rs := ms.db.Model(&s).Updates(ScheduleState{
				JobID: stat.JobID,
				Host:  stat.Host,
			}); rs.Error != nil {
				return rs.Error
			}
		} else {
			if rs := ms.db.Create(stat); rs.Error != nil {
				return rs.Error
			}
		}

		return nil
	})

	return err
}

func (ms *metaStore) DeleteScheduler(stat *ScheduleState) error {
	if rs := ms.db.Where(stat).Delete(&ScheduleState{}); rs.Error != nil {
		return rs.Error
	}
	return nil
}
