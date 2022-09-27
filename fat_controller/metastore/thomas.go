package metastore

import (
	"github.com/BabySid/proto/sodor"
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
