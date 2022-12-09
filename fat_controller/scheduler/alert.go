package scheduler

import (
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/structpb"
	"sodor/fat_controller/alert"
	"sodor/fat_controller/metastore"
	"time"
)

func giveAlertByGroupID(gid uint, param map[string]interface{}) error {
	in := sodor.AlertGroup{Id: int32(gid)}
	err := metastore.GetInstance().ShowAlertGroup(&in)
	if err != nil {
		return err
	}

	alertInsID := time.Now().Unix()
	alertsMap := in.PluginParams.AsMap()
	for name, preParam := range alertsMap {
		plugin := sodor.AlertPluginName(sodor.AlertPluginName_value[name])
		err = alert.GetInstance().GetAlertPlugin(plugin).GiveAlarm(param)
		status := "OK"
		if err != nil {
			status = err.Error()
		}

		preParamMap := preParam.(map[string]interface{})
		for k, v := range param {
			preParamMap[k] = v
		}

		pluginValue, err := structpb.NewStruct(preParamMap)
		if err != nil {
			log.Warnf("structpb.NewStruct(%+v) failed. err=%s", preParamMap, err)
			continue
		}
		ins := sodor.AlertGroupInstance{
			InstanceId:  int32(alertInsID),
			GroupId:     int32(gid),
			PluginName:  name,
			PluginValue: pluginValue,
			StatusMsg:   status,
		}
		err = metastore.GetInstance().InsertAlertGroupInstance(&ins)
		if err != nil {
			log.Warnf("InsertAlertGroupInstance failed. err=%s", err)
		}
	}
	return nil
}
