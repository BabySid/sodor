package alert

import (
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"sodor/fat_controller/metastore"
	"sync"
)

const (
	SystemAlertGroupName = "SODOR"
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

func (s *sodorAlert) AlertGroupID() int32 {
	return s.alertGroupID
}

func (s *sodorAlert) ResetAlertGroupID() error {
	ag, plugins, err := metastore.GetInstance().ShowSodorAlert(SystemAlertGroupName)
	if err != nil && err != metastore.ErrNotFound {
		return err
	}

	if ag == nil || plugins == nil {
		s.alertGroupID = 0
		s.alerts = nil
		log.Info("sodor alert is not set")
		return nil
	}

	s.alertGroupID = ag.Id

	s.alerts = make(map[int32]Alert)
	for id, plugin := range plugins.AlertPluginInstances {
		ding := NewDingDing(plugin.Dingding.Webhook, plugin.Dingding.Sign, plugin.Dingding.AtMobiles)
		s.alerts[int32(id)] = ding
	}

	log.Infof("sodor alert is set to group:%d plugins:%d", ag.Id, len(s.alerts))
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
