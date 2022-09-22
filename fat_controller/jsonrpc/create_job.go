package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) CreateJob(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.JobReply, *httpapi.JsonRpcError) {
	if params.Id != 0 {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("job.id must not be set"))
	}

	if err := checkTaskValid(params); err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams, httpapi.SysCodeMap[httpapi.InvalidParams], err)
	}

	exist, err := metastore.GetInstance().JobExist(params)
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	if exist {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("job exist"))
	}

	err = metastore.GetInstance().InsertJob(params)
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	ctx.ToLog("CreateJob Done: %+v", params)
	return &sodor.JobReply{Id: params.Id}, nil
}
