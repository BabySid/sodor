package grpc

import (
	"context"
	"github.com/BabySid/gorpc/grpc"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"sodor/fat_controller/scheduler"
)

func (s *Service) UpdateTaskInstance(ctx context.Context, task *sodor.TaskInstance) (*sodor.EmptyResponse, error) {
	ip, _ := grpc.GetPeerIPFromGRPC(ctx)
	log.WithFields(log.Fields{"job_instance_id": task.JobInstanceId, "job_id": task.JobId, "task_id": task.TaskId}).
		Infof("UpdateTaskInstance from %s", ip)

	_ = scheduler.GetInstance().UpdateTaskInstance(task)
	return &sodor.EmptyResponse{}, nil
}
