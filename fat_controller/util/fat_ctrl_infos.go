package util

import (
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/config"
)

func BuildFatCtrlInfos() *sodor.FatCtrlInfos {
	var infos sodor.FatCtrlInfos

	infos.FatCtrlInfos = make([]*sodor.FatCtrlInfo, 1)
	infos.FatCtrlInfos[0] = &sodor.FatCtrlInfo{
		Name:    config.GetInstance().AppName,
		Version: config.GetInstance().AppVersion,
		Proto:   "grpc",
		Host:    config.GetInstance().LocalIP,
		Port:    int32(config.GetInstance().Port),
	}

	return &infos
}
