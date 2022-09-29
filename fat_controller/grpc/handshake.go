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

func (s *Service) HandShake(ctx context.Context, req *sodor.ThomasInstance) (*sodor.ThomasReply, error) {
	ip, _ := grpc.GetPeerIPFromGRPC(ctx)
	log.Infof("HandShake from %s %d", ip, req.Id)

	if err := metastore.GetInstance().UpsertThomas(req); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &sodor.ThomasReply{Id: req.Id}, nil
}
