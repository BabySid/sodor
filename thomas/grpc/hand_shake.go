package grpc

import (
	"context"
	"github.com/BabySid/gorpc/grpc"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sodor/thomas/config"
	"sodor/thomas/fat_ctrl"
	"sodor/thomas/util"
)

func (s *Service) HandShake(ctx context.Context, req *sodor.HandShakeWithThomasRequest) (*sodor.ThomasInfo, error) {
	ip, _ := grpc.GetPeerIPFromGRPC(ctx)
	log.Infof("HandShake from %s. id=%d", ip, req.Thomas.Id)

	config.GetInstance().ThomasID = req.Thomas.Id
	_ = fat_ctrl.GetInstance().UpdateFatCtrlHost(req.FatCtrls)

	resp, err := util.BuildThomasInfo()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
