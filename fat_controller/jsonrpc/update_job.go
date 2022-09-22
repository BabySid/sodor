package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) UpdateJob(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.JobReply, *httpapi.JsonRpcError) {
	if params.Id == 0 {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("job.id must be set"))
	}

	if err := checkTaskValid(params); err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams, httpapi.SysCodeMap[httpapi.InvalidParams], err)
	}

	exist, err := metastore.GetInstance().JobExist(params)
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	if !exist {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("job not exist"))
	}

	err = metastore.GetInstance().UpdateJob(params)
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	ctx.ToLog("UpdateJob Done: %+v", params)
	return nil, nil
}
