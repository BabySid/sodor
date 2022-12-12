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
	"sodor/thomas/task_runner"
	"time"
)

func (s *Service) RunTask(ctx context.Context, task *sodor.RunTaskRequest) (*sodor.EmptyResponse, error) {
	ip, _ := grpc.GetPeerIPFromGRPC(ctx)
	log.Infof("RunTask from %s. task.id=%d task.insId=%d", ip, task.Task.Id, task.TaskInstance.Id)

	c, err := task_runner.GetTaskEnv().SetUp(task)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	go func() {
		_ = <-c.Start()
		resp, err := task_runner.GetTaskEnv().GetTaskResponse(c.ID)
		if resp == nil {
			log.Warnf("GetTaskResponse(%s) return nil resp. err = %v", c.ID, err)
			return
		}

		for {
			err = fat_ctrl.GetInstance().UpdateTaskInstance(resp)
			if err == nil {
				task_runner.GetTaskEnv().Remove(c.ID)
				break
			}

			log.Warnf("UpdateTaskInstance(%s) failed. retry after %v", c.ID, config.GetInstance().RetryInterval)
			time.Sleep(config.GetInstance().RetryInterval)
		}
	}()

	return &sodor.EmptyResponse{}, nil
}
