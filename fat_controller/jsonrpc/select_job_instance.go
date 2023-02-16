package jsonrpc

import (
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) SelectJobInstances(ctx api.Context, params *sodor.JobInstance) (*sodor.JobInstances, *api.JsonRpcError) {
	jobs, err := metastore.GetInstance().SelectJobInstance(params.JobId)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}
	ctx.Log("SelectJobInstances Done: %d", len(jobs.GetJobInstances()))
	return jobs, nil
}
