package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) SelectJob(ctx api.Context, params *sodor.Job) (*sodor.Job, *api.JsonRpcError) {
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

	ctx.Log("SelectJob Done: %+v", params)
	return &job, nil
}
