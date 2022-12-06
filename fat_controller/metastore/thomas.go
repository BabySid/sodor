package metastore

import (
	"errors"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sodor/fat_controller/config"
	"time"
)

var (
	NotFoundErr = errors.New("thomas not found")
)

func (ms *metaStore) UpsertThomas(req *sodor.ThomasInfo) error {
	var thomas Thomas
	_ = toThomas(req, &thomas)

	rs := ms.db.Transaction(func(tx *gorm.DB) error {
		gobase.True(thomas.ID > 0)
		rs := tx.Model(&thomas).Select(thomas.UpdateFields()).Updates(thomas)

		if rs.Error != nil {
			return rs.Error
		}

		if rs.RowsAffected == 0 {
			log.Warnf("unknown thomas. maybe it has been dropped. %+v", thomas)
			return NotFoundErr
		}

		var tIns ThomasInstance
		tIns.ThomasID = int32(thomas.ID)
		tIns.Metrics = thomas.Metrics
		if rs = tx.Create(&tIns); rs.Error != nil {
			return rs.Error
		}

		subQuery := tx.Model(&ThomasInstance{}).Select("id").Where("thomas_id = ?", tIns.ThomasID).
			Order("id desc").Offset(int(config.GetInstance().MaxThomasInstance)).Limit(1024)
		sub := tx.Table("(?) as tbl", subQuery).Select("tbl.id")
		if rs = tx.Where("id in (?)", sub).Delete(&ThomasInstance{}); rs.Error != nil {
			return rs.Error
		}

		return nil
	})

	return rs
}

func (ms *metaStore) UpdateThomasStatus(id int32, status string) error {
	var thomas Thomas
	thomas.ID = uint(id)

	if rs := ms.db.Model(&thomas).Updates(Thomas{Status: status}); rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (ms *metaStore) ThomasExist(host string, port int32) (bool, error) {
	id, err := ms.getThomasByHostPort(host, port)
	if err != nil {
		return false, err
	}

	if id > 0 {
		return true, nil
	}

	return false, nil
}

func (ms *metaStore) getThomasByHostPort(host string, port int32) (uint, error) {
	type apiThomas struct {
		TableModel
	}

	var t Thomas
	t.Host = host
	t.Port = int(port)

	var ts apiThomas
	rs := ms.db.Model(&Thomas{}).Where(&t).Limit(1).Find(&ts)
	if rs.Error != nil {
		return 0, rs.Error
	}

	return ts.ID, nil
}

func (ms *metaStore) ThomasExistByID(id int32) (bool, error) {
	type apiThomas struct {
		TableModel
	}

	var t Thomas
	t.ID = uint(id)

	var ts apiThomas
	rs := ms.db.Model(&Thomas{}).Where(&t).Limit(1).Find(&ts)
	if rs.Error != nil {
		return false, rs.Error
	}

	if ts.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (ms *metaStore) AddThomas(thomas *sodor.ThomasInfo) error {
	var out Thomas
	if err := toSimpleThomas(thomas, &out); err != nil {
		return err
	}

	rs := ms.db.Create(&out)
	if rs.Error != nil {
		return rs.Error
	}

	thomas.Id = int32(out.ID)

	return nil
}

func (ms *metaStore) ShowThomas(in *sodor.ThomasInfo, out *sodor.ThomasInstance) error {
	gobase.True(in.Id > 0)

	var thomas Thomas
	rs := ms.db.Limit(1).Find(&thomas, in.Id)
	if rs.Error != nil {
		return rs.Error
	}

	if thomas.ID == 0 {
		return ErrNotFound
	}

	out.Thomas = &sodor.ThomasInfo{}
	if err := fromThomas(&thomas, out.Thomas); err != nil {
		return err
	}

	var ins []*ThomasInstance
	if rs = ms.db.Where(&ThomasInstance{ThomasID: in.Id}).Find(&ins); rs.Error != nil {
		return rs.Error
	}

	out.Metrics = make([]*sodor.ThomasMetrics, len(ins))
	for i, t := range ins {
		var metrics sodor.ThomasMetrics
		if err := fromThomasMetrics(t, &metrics); err != nil {
			return err
		}
		out.Metrics[i] = &metrics
	}
	return nil
}

func (ms *metaStore) DropThomas(thomas *sodor.ThomasInfo) error {
	gobase.True(thomas.Id > 0)

	err := ms.db.Transaction(func(tx *gorm.DB) error {
		var t Thomas
		t.ID = uint(thomas.Id)

		if rs := tx.Delete(&t); rs.Error != nil {
			return rs.Error
		}

		if rs := tx.Where("thomas_id = ?", t.ID).Delete(&ThomasInstance{}); rs.Error != nil {
			return rs.Error
		}

		return nil
	})

	return err
}

func (ms *metaStore) SelectValidThomas(host string) (*Thomas, error) {
	var thomas Thomas
	thomas.Host = host
	if rs := ms.db.Scopes(filterValidThomas).Order("id desc").Where(thomas).Limit(1).Find(&thomas); rs.Error != nil {
		return nil, rs.Error
	}
	return &thomas, nil
}

func filterValidThomas(db *gorm.DB) *gorm.DB {
	return db.Where("heart_beat_time >= ?", time.Now().Unix()-maxThomasLife)
}

func (ms *metaStore) SelectInvalidThomas() ([]Thomas, error) {
	var thomas []Thomas
	if rs := ms.db.Scopes(filterInvalidThomas).Find(&thomas); rs.Error != nil {
		return nil, rs.Error
	}
	return thomas, nil
}

func filterInvalidThomas(db *gorm.DB) *gorm.DB {
	return db.Where("heart_beat_time < ?", time.Now().Unix()-maxThomasLife)
}

func (ms *metaStore) ListAllThomas() (*sodor.ThomasInfos, error) {
	var thomas []Thomas

	if rs := ms.db.Find(&thomas); rs.Error != nil {
		return nil, rs.Error
	}

	var all sodor.ThomasInfos
	all.ThomasInfos = make([]*sodor.ThomasInfo, len(thomas))

	for i, t := range thomas {
		var ti sodor.ThomasInfo
		err := fromThomas(&t, &ti)
		if err != nil {
			return nil, err
		}

		all.ThomasInfos[i] = &ti
	}

	return &all, nil
}
