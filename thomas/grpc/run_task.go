package grpc

import (
	"context"
	"github.com/BabySid/proto/sodor"
)

func (s *Service) RunTask(ctx context.Context, task *sodor.RunTaskRequest) (*sodor.EmptyResponse, error) {
	return &sodor.EmptyResponse{}, nil
}
