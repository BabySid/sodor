package grpc

import (
	"context"
	"github.com/BabySid/gorpc/util"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"sodor/fat_controller/scheduler"
)

func (s *Service) UpdateTaskInstance(ctx context.Context, task *sodor.TaskInstance) (*sodor.EmptyResponse, error) {
	ip, _ := util.GetPeerIPFromGRPC(ctx)
	log.WithFields(log.Fields{"job_instance_id": task.JobInstanceId, "job_id": task.JobId, "task_id": task.TaskId}).
		Infof("UpdateTaskInstance(taskInsID:%d taskID:%d) with exit_code:%d from %s", task.Id, task.TaskId, task.ExitCode, ip)

	_ = scheduler.GetInstance().UpdateTaskInstance(task)
	return &sodor.EmptyResponse{}, nil
}
