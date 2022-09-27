package grpc

import (
	"context"
	"github.com/BabySid/gorpc/grpc"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sodor/fat_controller/metastore"
)

func (s *Service) HandShakeWithThomas(ctx context.Context, req *sodor.ThomasHandShakeReq) (*sodor.ThomasHandShakeResp, error) {
	ip, _ := grpc.GetPeerIPFromGRPC(ctx)
	log.Infof("HandShakeWithThomas from %s %d", ip, req.Id)

	if err := metastore.GetInstance().UpsertThomas(req); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &sodor.ThomasHandShakeResp{Id: req.Id}, nil
}
