package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
)

func (s *Service) AddThomas(ctx *httpapi.APIContext, params *sodor.ThomasInfo) (*sodor.ThomasReply, *httpapi.JsonRpcError) {
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

	params.ThomasType = sodor.ThomasType_Thomas_Static
	if err = metastore.GetInstance().AddThomas(params); err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	err = scheduler.PingThomas(params.Id, params.Host, int(params.Port))

	ctx.ToLog("AddThomas Done: %+v", params)

	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError,
			httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	return &sodor.ThomasReply{Id: params.Id}, nil
}