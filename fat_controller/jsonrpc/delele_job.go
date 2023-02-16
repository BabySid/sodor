package jsonrpc

import (
	"errors"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
)

func (s *Service) DeleteJob(ctx api.Context, params *sodor.Job) (*sodor.JobReply, *api.JsonRpcError) {
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

	if err = metastore.GetInstance().DeleteJob(params); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	err = scheduler.GetInstance().Remove(params)
	gobase.True(err == nil)

	ctx.Log("DeleteJob Done: %+v", params)
	return &sodor.JobReply{Id: params.Id}, nil
}
