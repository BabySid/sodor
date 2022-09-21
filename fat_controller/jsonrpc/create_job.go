package jsonrpc

import (
	"fmt"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) CreateJob(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.JobReply, *httpapi.JsonRpcError) {
	if params.Id != 0 {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams, "params.id cannot be set", nil)
	}

	if len(params.Name) >= metastore.MaxNameLen {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			fmt.Sprintf("params.name is long than %d", metastore.MaxNameLen), nil)
	}

	if err := checkTaskValid(params); err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams, err.Error(), nil)
	}

	err := metastore.GetInstance().InsertJob(params)
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, err.Error(), nil)
	}

	ctx.ToLog("CreateJob Done: %+v", params)
	return &sodor.JobReply{Id: params.Id}, nil
}
