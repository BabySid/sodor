package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ShowThomas(ctx *httpapi.APIContext, params *sodor.ThomasInfo) (*sodor.ThomasInstance, *httpapi.JsonRpcError) {
	if params.Id == 0 {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("invalid params of thomas"))
	}

	exist, err := metastore.GetInstance().ThomasExistByID(params.Id)
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	if !exist {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("thomas not exist"))
	}

	var out sodor.ThomasInstance
	if err = metastore.GetInstance().ShowThomas(params, &out); err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	ctx.ToLog("ShowThomas Done: %+v", params)
	return &out, nil
}
