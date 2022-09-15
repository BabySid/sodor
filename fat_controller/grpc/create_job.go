package grpc

import (
	"context"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
)

func (s *Service) CreateJob(ctx context.Context, in *sodor.Job) (*sodor.Reply, error) {
	log.Infof("job %+v", in)
	return &sodor.Reply{}, nil
}
