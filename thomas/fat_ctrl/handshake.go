package fat_ctrl

import (
	"context"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"os"
	"sodor/fat_controller/config"
)

func (fc *FatCtrl) HandShake() {
	if fc.thomasID == 0 || len(fc.fatCtrlHosts) == 0 {
		log.Warn("wait for get fat_ctrl info. thomasID=%d fc.fatCtrlHosts.size=%d", fc.thomasID, len(fc.fatCtrlHosts))
		return
	}

	conn, err := fc.getFatCtrlHost()
	if err != nil {
		log.Warnf("getFatCtrlHost failed. err=%s", err)
		return
	}

	defer conn.Close()

	cli := sodor.NewFatControllerClient(conn)

	var req sodor.ThomasInfo
	req.Id = fc.thomasID
	req.Version = config.GetInstance().AppVersion
	req.Name = config.GetInstance().AppName
	req.Proto = "grpc"
	req.Host = config.GetInstance().LocalIP
	req.Port = int32(config.GetInstance().Port)
	req.Pid = int32(os.Getpid())
	req.StartTime = fc.startTime
	m := fc.getMetrics()
	metrics, err := structpb.NewStruct(m)
	if err != nil {
		log.Warnf("structpb.NewStruct failed. raw_data=%v err=%s", m, err)
		return
	}
	req.LatestMetrics = metrics

	_, err = cli.HandShake(context.Background(), &req)
	if s, ok := status.FromError(err); ok {
		if s != nil {
			log.Warnf("HandShake to thomas failed. code=%d, msg=%s", s.Code(), s.Message())
		}
	} else {
		if err != nil {
			log.Warnf("HandShake to thomas failed. err=%s", err)
		}
	}
}

const (
	cpuUsage        = "cpu_used_percent"
	memUsage        = "mem_used_percent"
	diskUsagePrefix = "disk_used_percent_"
)

func (fc *FatCtrl) getMetrics() map[string]interface{} {
	rs := make(map[string]interface{})
	rs[cpuUsage] = gobase.GetCPUUsage()
	rs[memUsage] = gobase.GetMEMUsage()
	dp := gobase.GetDiskPartitionUsedPercent()
	for d, up := range dp {
		rs[diskUsagePrefix+d] = up
	}
	return rs
}
