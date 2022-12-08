package alert

import (
	"github.com/BabySid/proto/sodor"
	"sync"
)

type alert interface {
	GetName() string
	GetParams() []sodor.AlertPluginParams
	GiveAlarm(param map[string]interface{}) error
}

type alertFactory struct {
	alerts map[sodor.AlertPluginName]alert
}

func (af *alertFactory) GetAlertPlugins() []alert {
	rs := make([]alert, 0)

	for _, v := range af.alerts {
		rs = append(rs, v)
	}

	return rs
}

func (af *alertFactory) GetAlertPlugin(name sodor.AlertPluginName) alert {
	if a, ok := af.alerts[name]; ok {
		return a
	}

	return nil
}

var (
	once      sync.Once
	singleton *alertFactory
)

func GetInstance() *alertFactory {
	once.Do(func() {
		singleton = &alertFactory{
			alerts: map[sodor.AlertPluginName]alert{
				sodor.AlertPluginName_APN_DingDing: &dingDing{},
			},
		}
	})
	return singleton
}
