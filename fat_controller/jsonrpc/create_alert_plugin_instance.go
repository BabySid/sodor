package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) CreateAlertPluginInstance(ctx api.Context, params *sodor.AlertPluginInstance) (*sodor.AlertPluginReply, *api.JsonRpcError) {
	if err := checkAlertPluginValid(params, true); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, err)
	}

	exist, err := metastore.GetInstance().AlertPluginInstanceExist(params)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if exist {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("alert_plugin exist"))
	}

	if err = metastore.GetInstance().InsertAlertPluginInstance(params); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	ctx.Log("CreateAlertPluginInstance Done: %+v", params)
	return &sodor.AlertPluginReply{Id: params.Id}, nil
}
