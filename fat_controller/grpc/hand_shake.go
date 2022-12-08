package grpc

import (
	"context"
	"github.com/BabySid/gorpc/grpc"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/util"
)

func (s *Service) HandShake(ctx context.Context, req *sodor.ThomasInfo) (*sodor.FatCtrlInfos, error) {
	ip, _ := grpc.GetPeerIPFromGRPC(ctx)
	log.Infof("HandShake from %s req.Id=%d", ip, req.Id)

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid req.Id")
	}

	if err := metastore.GetInstance().UpsertThomas(req); err != nil {
		if err == metastore.NotFoundErr {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	infos := util.BuildFatCtrlInfos()
	return infos, nil
}
