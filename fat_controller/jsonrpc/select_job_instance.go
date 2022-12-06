package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) SelectJobInstances(ctx *httpapi.APIContext, params *sodor.JobInstance) (*sodor.JobTaskInstances, *httpapi.JsonRpcError) {
	jobs, err := metastore.GetInstance().SelectInstanceByJobID(params.JobId)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}
	ctx.ToLog("SelectJobInstances Done: %d %d", len(jobs.GetJobInstances().GetJobInstances()), len(jobs.GetTaskInstances().GetTaskInstances()))
	return jobs, nil
}
