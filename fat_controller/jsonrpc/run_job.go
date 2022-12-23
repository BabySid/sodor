package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
)

func (s *Service) RunJob(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.JobReply, *httpapi.JsonRpcError) {
	if params.Id == 0 {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("job.id must be set"))
	}

	exist, err := metastore.GetInstance().JobExist(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	if !exist {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("job not exist"))
	}

	var job sodor.Job
	job.Id = params.Id

	err = metastore.GetInstance().SelectJob(&job)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	_ = scheduler.GetInstance().RunJob(&job)

	ctx.ToLog("RunJob Done: %+v", params)
	return nil, nil
}
