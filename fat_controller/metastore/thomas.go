package metastore

import (
	"github.com/BabySid/proto/sodor"
	"gorm.io/gorm"
	"time"
)

func (ms *metaStore) UpsertThomas(job *sodor.ThomasHandShakeReq) error {
	var thomas Thomas
	_ = toThomas(job, &thomas)
	if thomas.ID > 0 {
		if rs := ms.db.Model(&thomas).Update("last_heartbeat_time", time.Now().Unix()); rs.Error != nil {
			return rs.Error
		}
	} else {
		if rs := ms.db.Create(&thomas); rs.Error != nil {
			return rs.Error
		}

		job.Id = int32(thomas.ID)
	}
	return nil
}

func (ms *metaStore) SelectThomas(host string) (*Thomas, error) {
	var thomas Thomas
	thomas.Host = host
	//ms.db.
	if rs := ms.db.Scopes(filterValidThomas).Last(&thomas); rs.Error != nil {
		return nil, rs.Error
	}
	return &thomas, nil
}

func filterValidThomas(db *gorm.DB) *gorm.DB {
	const MaxThomasLife = 300
	return db.Where("heart_beat_time >= ?", time.Now().Unix()-MaxThomasLife)
}
