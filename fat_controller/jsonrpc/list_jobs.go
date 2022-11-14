package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ListJobs(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.Jobs, *httpapi.JsonRpcError) {
	jobs, err := metastore.GetInstance().ListJobs()
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}
	ctx.ToLog("ListJobs Done: %d", len(jobs.GetJobs()))
	return jobs, nil
}
