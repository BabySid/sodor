package jsonrpc

import (
	"errors"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
)

func (s *Service) DeleteJob(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.JobReply, *httpapi.JsonRpcError) {
	if params.Id == 0 {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("job.id must be set"))
	}

	exist, err := metastore.GetInstance().JobExist(params)
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	if !exist {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("job not exist"))
	}

	if err = metastore.GetInstance().DeleteJob(params); err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	err = scheduler.GetInstance().Remove(params)
	gobase.True(err == nil)

	ctx.ToLog("DeleteJob Done: %+v", params)
	return &sodor.JobReply{Id: params.Id}, nil
}
