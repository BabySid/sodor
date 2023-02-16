package jsonrpc

import (
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) SelectTaskInstances(ctx api.Context, params *sodor.TaskInstance) (*sodor.TaskInstances, *api.JsonRpcError) {
	tasks, err := metastore.GetInstance().SelectTaskInstance(params.JobId, params.JobInstanceId)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}
	ctx.Log("SelectTaskInstances Done: %d", len(tasks.GetTaskInstances()))
	return tasks, nil
}
