package metastore

import (
	"github.com/BabySid/proto/sodor"
	"gorm.io/gorm"
	"time"
)

func (ms *metaStore) UpsertThomas(req *sodor.ThomasInstance) error {
	var thomas Thomas
	_ = toThomas(req, &thomas)

	if thomas.ID == 0 {
		id, err := ms.getThomasByHostPort(req.Host, req.Port)
		if err != nil {
			return err
		}
		thomas.ID = id
	}

	if thomas.ID > 0 {
		if rs := ms.db.Save([]*Thomas{&thomas}); rs.Error != nil {
			return rs.Error
		}
	} else {
		if rs := ms.db.Create(&thomas); rs.Error != nil {
			return rs.Error
		}

		req.Id = int32(thomas.ID)
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

	var ts []apiThomas
	rs := ms.db.Model(&Thomas{}).Where(&t).Find(&ts)
	if rs.Error != nil {
		return 0, rs.Error
	}
	if rs.RowsAffected > 0 {
		return ts[0].ID, nil
	}

	return 0, nil
}

func (ms *metaStore) ThomasExistByID(id int32) (bool, error) {
	type apiThomas struct {
		TableModel
	}

	var t Thomas
	t.ID = uint(id)

	var ts []apiThomas
	rs := ms.db.Model(&Thomas{}).Where(&t).Find(&ts)
	if rs.Error != nil {
		return false, rs.Error
	}

	if rs.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func (ms *metaStore) AddThomas(thomas *sodor.ThomasInstance) error {
	var out Thomas
	if err := toSimpleThomas(thomas, &out); err != nil {
		return err
	}

	rs := ms.db.Create(&thomas)
	if rs.Error != nil {
		return rs.Error
	}

	thomas.Id = int32(out.ID)

	return nil
}

func (ms *metaStore) DropThomas(thomas *sodor.ThomasInstance) error {
	var out Thomas
	if err := toSimpleThomas(thomas, &out); err != nil {
		return err
	}

	rs := ms.db.Delete(&thomas)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (ms *metaStore) SelectValidThomas(host string) (*Thomas, error) {
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

func (ms *metaStore) ListAllThomas() (*sodor.ThomasInstances, error) {
	var thomas []Thomas

	if rs := ms.db.Find(&thomas); rs.Error != nil {
		return nil, rs.Error
	}

	var all sodor.ThomasInstances
	all.ThomasInstances = make([]*sodor.ThomasInstance, len(thomas))

	for i, t := range thomas {
		var ti sodor.ThomasInstance
		err := fromThomas(&t, &ti)
		if err != nil {
			return nil, err
		}

		all.ThomasInstances[i] = &ti
	}

	return nil, nil
}
