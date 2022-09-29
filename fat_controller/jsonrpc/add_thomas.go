package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) AddThomas(ctx *httpapi.APIContext, params *sodor.ThomasInstance) (*sodor.ThomasReply, *httpapi.JsonRpcError) {
	if params.Id != 0 || params.Pid != 0 || params.StartTime != 0 {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("invalid params of thomas"))
	}

	exist, err := metastore.GetInstance().ThomasExist(params.Host, params.Port)
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	if exist {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("thomas exist"))
	}

	if err = metastore.GetInstance().AddThomas(params); err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	ctx.ToLog("AddThomas Done: %+v", params)
	return &sodor.ThomasReply{Id: params.Id}, nil
}
