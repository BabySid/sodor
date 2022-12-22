package alert

import (
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"sodor/fat_controller/metastore"
	"sync"
)

const (
	systemAlertGroupName = "SODOR"
)

type sodorAlert struct {
	alertGroupID int32
	alerts       map[int32]Alert
}

var (
	once      sync.Once
	singleton *sodorAlert
)

func GetInstance() *sodorAlert {
	once.Do(func() {
		singleton = &sodorAlert{}
		if err := singleton.ResetAlertGroupID(); err != nil {
			log.Fatalf("ResetAlertGroupID failed. err=%s", err)
		}
	})
	return singleton
}

func (s *sodorAlert) ResetAlertGroupID() error {
	ag, plugins, err := metastore.GetInstance().ShowSodorAlert(systemAlertGroupName)
	if err != nil && err != metastore.ErrNotFound {
		return err
	}

	if ag == nil || plugins == nil {
		s.alertGroupID = 0
		s.alerts = nil
		log.Info("sodor alert is not set")
		return nil
	}

	for id, plugin := range plugins.AlertPluginInstances {
		param := plugin.Plugin.(*sodor.AlertPluginInstance_Dingding)
		ding := NewDingDing(param.Dingding.Webhook, param.Dingding.Sign, param.Dingding.AtMobiles)
		s.alerts[int32(id)] = ding
	}

	log.Info("sodor alert is set to group:%d", ag.Id)
	return nil
}

func (s *sodorAlert) GiveAlert(msg string) {
	if len(s.alerts) == 0 {
		return
	}

	for id, v := range s.alerts {
		err := v.GiveAlarm(msg)
		status := "OK"
		if err != nil {
			status = err.Error()
		}
		his := sodor.AlertPluginInstanceHistory{
			InstanceId: id,
			GroupId:    s.alertGroupID,
			AlertMsg:   msg,
			StatusMsg:  status,
		}
		err = metastore.GetInstance().InsertAlertPluginInstanceHistory(&his)
		if err != nil {
			log.Warnf("SodorAlert failed. pluginInstanceID=%d err=%s", id, err)
		}
	}
}
