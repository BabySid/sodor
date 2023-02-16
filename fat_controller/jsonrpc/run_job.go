package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
)

func (s *Service) RunJob(ctx api.Context, params *sodor.Job) (*sodor.JobReply, *api.JsonRpcError) {
	if params.Id == 0 {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("job.id must be set"))
	}

	exist, err := metastore.GetInstance().JobExist(params)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if !exist {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("job not exist"))
	}

	var job sodor.Job
	job.Id = params.Id

	err = metastore.GetInstance().SelectJob(&job)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	_ = scheduler.GetInstance().RunJob(&job)

	ctx.Log("RunJob Done: %+v", params)
	return nil, nil
}
