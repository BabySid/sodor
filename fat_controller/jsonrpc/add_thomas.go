package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
)

func (s *Service) AddThomas(ctx api.Context, params *sodor.ThomasInfo) (*sodor.ThomasReply, *api.JsonRpcError) {
	if params.Id != 0 || params.Pid != 0 || params.StartTime != 0 {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("invalid params of thomas"))
	}

	exist, err := metastore.GetInstance().ThomasExist(params.Host, params.Port)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if exist {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("thomas exist"))
	}

	params.ThomasType = sodor.ThomasType_Thomas_Static
	if err = metastore.GetInstance().AddThomas(params); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	err = scheduler.PingThomas(params.Id, params.Host, int(params.Port))

	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	ctx.Log("AddThomas Done: %+v", params)
	return &sodor.ThomasReply{Id: params.Id}, nil
}
