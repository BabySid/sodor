package jsonrpc

import (
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ListJobs(ctx api.Context, _ *interface{}) (*sodor.Jobs, *api.JsonRpcError) {
	jobs, err := metastore.GetInstance().ListJobs()
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}
	ctx.Log("ListJobs Done: %d", len(jobs.GetJobs()))
	return jobs, nil
}
