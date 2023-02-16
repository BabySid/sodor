package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) UpdateAlertPluginInstance(ctx api.Context, params *sodor.AlertPluginInstance) (*sodor.AlertPluginReply, *api.JsonRpcError) {
	if err := checkAlertPluginValid(params, false); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, err)
	}

	exist, err := metastore.GetInstance().AlertPluginInstanceExist(params)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if !exist {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("alert_plugin_instance not exist"))
	}

	if err = metastore.GetInstance().UpdateAlertPluginInstance(params); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	ctx.Log("UpdateAlertPluginInstance Done: %+v", params)
	return &sodor.AlertPluginReply{Id: params.Id}, nil
}
