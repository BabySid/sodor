package grpc

import (
	"context"
	"github.com/BabySid/proto/sodor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sodor/thomas/task_runner"
)

func (s *Service) RunTask(ctx context.Context, task *sodor.RunTaskRequest) (*sodor.EmptyResponse, error) {
	env := &task_runner.TaskEnv{}
	c, err := env.SetUp(task)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	go func() {
		_ = <-c.Start()
	}()

	return &sodor.EmptyResponse{}, nil
}
