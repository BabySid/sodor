package metastore

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"gorm.io/gorm"
	"sodor/fat_controller/config"
)

func (ms *metaStore) AlertPluginInstanceExist(ap *sodor.AlertPluginInstance) (bool, error) {
	type apiAlertPlugin struct {
		gobase.TableModel
	}

	var a AlertPluginInstance
	if ap.Id > 0 {
		a.ID = uint(ap.Id)
	} else if ap.Name != "" {
		a.Name = ap.Name
	}

	var alert apiAlertPlugin
	rs := ms.db.Model(&AlertPluginInstance{}).Where(&a).Limit(1).Find(&alert)
	if rs.Error != nil {
		return false, rs.Error
	}

	if alert.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (ms *metaStore) AlertPluginInstanceUsedInAlertGroup(ap *sodor.AlertPluginInstance) (bool, error) {
	var ag AlertGroup
	rs := ms.db.Where(fmt.Sprintf("json_contains(plugin_instance, '%d')", ap.Id)).Limit(1).Find(&ag)
	if rs.Error != nil {
		return false, rs.Error
	}

	if ag.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (ms *metaStore) AlertGroupExist(ag *sodor.AlertGroup) (bool, error) {
	type apiAlert struct {
		gobase.TableModel
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

func (ms *metaStore) InsertAlertPluginInstance(plugin *sodor.AlertPluginInstance) error {
	gobase.True(plugin.Id == 0)

	var out AlertPluginInstance
	if err := toAlertPluginInstance(plugin, &out); err != nil {
		return err
	}

	rs := ms.db.Create(&out)
	if rs.Error != nil {
		return rs.Error
	}

	plugin.Id = int32(out.ID)
	return nil
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

		if rs := tx.Where(AlertPluginInstanceHistory{GroupID: int32(t.ID)}).Delete(&AlertPluginInstanceHistory{}); rs.Error != nil {
			return rs.Error
		}

		return nil
	})
	return err
}

func (ms *metaStore) DeleteAlertPluginInstance(ins *sodor.AlertPluginInstance) error {
	gobase.True(ins.Id > 0)

	err := ms.db.Transaction(func(tx *gorm.DB) error {
		var t AlertPluginInstance
		t.ID = uint(ins.Id)

		if rs := tx.Delete(&t); rs.Error != nil {
			return rs.Error
		}

		if rs := tx.Where(AlertPluginInstanceHistory{InstanceId: int32(t.ID)}).Delete(&AlertPluginInstanceHistory{}); rs.Error != nil {
			return rs.Error
		}

		return nil
	})

	return err
}

func (ms *metaStore) UpdateAlertPluginInstance(in *sodor.AlertPluginInstance) error {
	gobase.True(in.Id > 0)

	var out AlertPluginInstance
	if err := toAlertPluginInstance(in, &out); err != nil {
		return err
	}

	if rst := ms.db.Model(&out).Select(out.UpdateFields()).Updates(out); rst.Error != nil {
		return rst.Error
	}

	return nil
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

func (ms *metaStore) ListAlertPluginInstances(instance ...*sodor.AlertPluginInstance) (*sodor.AlertPluginInstances, error) {
	var aps []AlertPluginInstance

	var rs *gorm.DB
	if len(instance) == 0 {
		rs = ms.db.Find(&aps)
	} else {
		pk := make([]uint, len(instance))
		for i, v := range instance {
			pk[i] = uint(v.Id)
		}

		rs = ms.db.Find(&aps, pk)
	}
	if rs.Error != nil {
		return nil, rs.Error
	}

	var all sodor.AlertPluginInstances
	all.AlertPluginInstances = make([]*sodor.AlertPluginInstance, len(aps))

	for i, t := range aps {
		var ap sodor.AlertPluginInstance
		err := fromAlertPluginInstance(&t, &ap)
		if err != nil {
			return nil, err
		}
		all.AlertPluginInstances[i] = &ap
	}

	return &all, nil
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

func (ms *metaStore) ShowSodorAlert(name string) (*sodor.AlertGroup, *sodor.AlertPluginInstances, error) {
	var ag AlertGroup
	rs := ms.db.Model(&AlertGroup{}).Where(AlertGroup{Name: name}).Limit(1).Find(&ag)
	if rs.Error != nil {
		return nil, nil, rs.Error
	}

	if rs.RowsAffected == 0 {
		return nil, nil, ErrNotFound
	}

	var sag sodor.AlertGroup
	sag.Id = int32(ag.ID)
	return ms.ShowAlertGroup(&sag)
}

func (ms *metaStore) ShowAlertGroup(group *sodor.AlertGroup) (*sodor.AlertGroup, *sodor.AlertPluginInstances, error) {
	gobase.True(group.Id > 0)

	var ag AlertGroup
	rs := ms.db.Limit(1).Find(&ag, group.Id)
	if rs.Error != nil {
		return nil, nil, rs.Error
	}

	if rs.RowsAffected == 0 {
		return nil, nil, ErrNotFound
	}

	if err := fromAlertGroup(&ag, group); err != nil {
		return nil, nil, err
	}

	pluginIns := make([]*sodor.AlertPluginInstance, len(ag.PluginInstance))
	for i, v := range ag.PluginInstance {
		pluginIns[i] = &sodor.AlertPluginInstance{Id: int32(v)}
	}

	out, err := ms.ListAlertPluginInstances(pluginIns...)
	if err != nil {
		return nil, nil, err
	}

	return group, out, nil
}

func (ms *metaStore) InsertAlertPluginInstanceHistory(his *sodor.AlertPluginInstanceHistory) error {
	gobase.True(his.Id == 0)

	err := ms.db.Transaction(func(tx *gorm.DB) error {
		var out AlertPluginInstanceHistory
		toAlertPluginInstanceHistory(his, &out)

		rs := tx.Create(&out)
		if rs.Error != nil {
			return rs.Error
		}

		his.Id = int32(out.ID)

		if config.GetInstance().MaxAlertHistory <= 0 {
			return nil
		}

		var hisIds []uint
		rs = tx.Model(&AlertPluginInstanceHistory{}).Where(AlertPluginInstanceHistory{GroupID: his.GroupId, InstanceId: his.InstanceId}).
			Order("id desc").Offset(int(config.GetInstance().MaxAlertHistory)).Limit(1024).
			Pluck("id", &hisIds)
		if rs.Error != nil {
			return rs.Error
		}

		if len(hisIds) == 0 {
			return nil
		}

		if rs = tx.Where("id in (?)", hisIds).Delete(&AlertPluginInstanceHistory{}); rs.Error != nil {
			return rs.Error
		}

		return nil
	})

	return err
}

func (ms *metaStore) ShowAlertPluginInstanceHistories(his *sodor.AlertPluginInstanceHistory) (*sodor.AlertPluginInstanceHistories, error) {
	var histories []*AlertPluginInstanceHistory
	rs := ms.db.Model(&AlertPluginInstanceHistory{}).Where(&AlertPluginInstanceHistory{GroupID: his.GroupId, InstanceId: his.InstanceId}).Find(&histories)
	if rs.Error != nil {
		return nil, rs.Error
	}

	var apih sodor.AlertPluginInstanceHistories
	apih.AlertPluginInstanceHistory = make([]*sodor.AlertPluginInstanceHistory, len(histories))

	for i, ag := range histories {
		var ins sodor.AlertPluginInstanceHistory
		if err := fromAlertPluginInstanceHistory(ag, &ins); err != nil {
			return nil, err
		}
		apih.AlertPluginInstanceHistory[i] = &ins
	}

	return &apih, nil
}
