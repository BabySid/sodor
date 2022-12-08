package util

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/structpb"
	"sodor/thomas/config"
)

func BuildThomasInfo() (*sodor.ThomasInfo, error) {
	var req sodor.ThomasInfo
	req.Id = config.GetInstance().ThomasID
	req.Version = config.GetInstance().AppVersion
	req.Name = config.GetInstance().AppName
	req.Proto = "grpc"
	req.Host = config.GetInstance().LocalIP
	req.Port = int32(config.GetInstance().Port)
	req.Pid = config.GetInstance().Pid
	req.StartTime = config.GetInstance().StartTime
	m := getMetrics()
	metrics, err := structpb.NewStruct(m)
	if err != nil {
		log.Warnf("structpb.NewStruct failed. raw_data=%v err=%s", m, err)
		return nil, err
	}
	req.LatestMetrics = metrics

	return &req, nil
}

const (
	cpuUsage        = "cpu_used_percent"
	memUsage        = "mem_used_percent"
	diskUsagePrefix = "disk_used_percent_"
)

func getMetrics() map[string]interface{} {
	rs := make(map[string]interface{})
	rs[cpuUsage] = gobase.GetCPUUsage()
	rs[memUsage] = gobase.GetMEMUsage()
	dp := gobase.GetDiskPartitionUsedPercent()
	for d, up := range dp {
		rs[diskUsagePrefix+d] = up
	}
	return rs
}
