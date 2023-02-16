package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) DropThomas(ctx api.Context, params *sodor.ThomasInfo) (*sodor.ThomasReply, *api.JsonRpcError) {
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

	if err = metastore.GetInstance().DropThomas(params); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}
	ctx.Log("DropThomas Done: %+v", params)
	return &sodor.ThomasReply{Id: params.Id}, nil
}
