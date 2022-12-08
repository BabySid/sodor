package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) DropThomas(ctx *httpapi.APIContext, params *sodor.ThomasInfo) (*sodor.ThomasInfo, *httpapi.JsonRpcError) {
	if params.Id == 0 {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("invalid params of thomas"))
	}

	exist, err := metastore.GetInstance().ThomasExistByID(params.Id)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	if !exist {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("thomas not exist"))
	}

	if err = metastore.GetInstance().DropThomas(params); err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}
	ctx.ToLog("DropThomas Done: %+v", params)
	return &sodor.ThomasInfo{Id: params.Id}, nil
}
