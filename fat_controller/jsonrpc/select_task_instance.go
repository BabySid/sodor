package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) SelectTaskInstances(ctx *httpapi.APIContext, params *sodor.TaskInstance) (*sodor.TaskInstances, *httpapi.JsonRpcError) {
	tasks, err := metastore.GetInstance().SelectTaskInstance(params.JobId, params.JobInstanceId)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}
	ctx.ToLog("SelectTaskInstances Done: %d", len(tasks.GetTaskInstances()))
	return tasks, nil
}
