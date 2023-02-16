package fat_ctrl

import (
	"context"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sodor/thomas/config"
	"sodor/thomas/util"
)

func (fc *FatCtrl) HandShake() {
	rpcCli, err := fc.getFatCtrlConn()
	if err != nil {
		log.Warnf("getFatCtrlConn failed. err=%s", err)
		return
	}

	defer rpcCli.Close()

	cli := sodor.NewFatControllerClient(rpcCli.UnderlyingHandle().(*grpc.ClientConn))

	req, err := util.BuildThomasInfo()
	if err != nil {
		log.Warnf("BuildThomasInfo failed. err=%s", err)
		return
	}

	resp, err := cli.HandShake(context.Background(), req)
	if s, ok := status.FromError(err); ok {
		if s != nil {
			if s.Code() == codes.NotFound {
				// reset the thomasID because the thomasID has been dropped
				config.GetInstance().ThomasID = 0
			}
			log.Warnf("HandShake to fat_ctrl failed. code=%d, msg=%s", s.Code(), s.Message())
			return
		}
	} else {
		if err != nil {
			log.Warnf("HandShake to fat_ctrl failed. err=%s", err)
			return
		}
	}

	_ = fc.UpdateFatCtrlHost(resp)
}
