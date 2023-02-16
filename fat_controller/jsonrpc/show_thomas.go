package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ShowThomas(ctx api.Context, params *sodor.ThomasInfo) (*sodor.ThomasInstance, *api.JsonRpcError) {
	if params.Id == 0 {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("invalid params of thomas"))
	}

	exist, err := metastore.GetInstance().ThomasExistByID(params.Id)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if !exist {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("thomas not exist"))
	}

	var out sodor.ThomasInstance
	if err = metastore.GetInstance().ShowThomas(params, &out); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	ctx.Log("ShowThomas Done: %+v", params)
	return &out, nil
}
