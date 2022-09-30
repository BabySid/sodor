package grpc

import (
	"context"
	"github.com/BabySid/gorpc/grpc"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
)

func (s *Service) HandShake(ctx context.Context, req *sodor.FatCtrlInfo) (*sodor.FatCtrlReply, error) {
	ip, _ := grpc.GetPeerIPFromGRPC(ctx)
	log.Infof("HandShake from %s %d", ip, req.Id)

	_ = s.updateFatCtrlHost(req.Host, int(req.Port))

	return &sodor.FatCtrlReply{Id: req.Id}, nil
}
