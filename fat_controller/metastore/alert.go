package metastore

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"gorm.io/gorm"
)

func (ms *metaStore) AlertGroupExist(ag *sodor.AlertGroup) (bool, error) {
	type apiAlert struct {
		TableModel
	}

	var a AlertGroup
	if ag.Id > 0 {
		a.ID = uint(ag.Id)
	} else if ag.Name != "" {
		a.Name = ag.Name
	}

	var alert apiAlert
	rs := ms.db.Model(&AlertGroup{}).Where(&a).Limit(1).Find(&alert)
	if rs.Error != nil {
		return false, rs.Error
	}

	if alert.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (ms *metaStore) InsertAlertGroup(alert *sodor.AlertGroup) error {
	gobase.True(alert.Id == 0)

	var out AlertGroup
	if err := toAlertGroup(alert, &out); err != nil {
		return err
	}

	rs := ms.db.Create(&out)
	if rs.Error != nil {
		return rs.Error
	}

	alert.Id = int32(out.ID)
	return nil
}

func (ms *metaStore) DeleteAlertGroup(alert *sodor.AlertGroup) error {
	gobase.True(alert.Id > 0)

	err := ms.db.Transaction(func(tx *gorm.DB) error {
		var t AlertGroup
		t.ID = uint(alert.Id)

		if rs := tx.Delete(&t); rs.Error != nil {
			return rs.Error
		}

		if rs := tx.Where(AlertGroupInstance{GroupID: int32(t.ID)}).Delete(&AlertGroupInstance{}); rs.Error != nil {
			return rs.Error
		}

		return nil
	})
	return err
}

func (ms *metaStore) UpdateAlertGroup(alert *sodor.AlertGroup) error {
	gobase.True(alert.Id > 0)

	var out AlertGroup
	if err := toAlertGroup(alert, &out); err != nil {
		return err
	}

	if rst := ms.db.Model(&out).Select(out.UpdateFields()).Updates(out); rst.Error != nil {
		return rst.Error
	}

	return nil
}

func (ms *metaStore) ListAlertGroups() (*sodor.AlertGroups, error) {
	var ags []AlertGroup

	if rs := ms.db.Find(&ags); rs.Error != nil {
		return nil, rs.Error
	}

	var all sodor.AlertGroups
	all.AlertGroups = make([]*sodor.AlertGroup, len(ags))

	for i, t := range ags {
		var ag sodor.AlertGroup
		err := fromAlertGroup(&t, &ag)
		if err != nil {
			return nil, err
		}

		all.AlertGroups[i] = &ag
	}

	return &all, nil
}

func (ms *metaStore) ShowAlertGroup(in *sodor.AlertGroup) error {
	gobase.True(in.Id > 0)

	var ag AlertGroup
	rs := ms.db.Limit(1).Find(&ag, in.Id)
	if rs.Error != nil {
		return rs.Error
	}

	if rs.RowsAffected == 0 {
		return ErrNotFound
	}

	if err := fromAlertGroup(&ag, in); err != nil {
		return err
	}

	return nil
}

func (ms *metaStore) ShowAlertGroupInstanceByGroupID(gID int32) (*sodor.AlertGroupInstances, error) {
	var agIns []*AlertGroupInstance
	rs := ms.db.Model(&AlertGroupInstance{}).Where(&AlertGroupInstance{GroupID: gID}).Find(&agIns)
	if rs.Error != nil {
		return nil, rs.Error
	}

	var ags sodor.AlertGroupInstances
	ags.AlertGroupInstances = make([]*sodor.AlertGroupInstance, len(agIns))

	for i, ag := range agIns {
		var ins sodor.AlertGroupInstance
		if err := fromAlertGroupInstance(ag, &ins); err != nil {
			return nil, err
		}
		ags.AlertGroupInstances[i] = &ins
	}

	return &ags, nil
}
