package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) DeleteAlertPluginInstance(ctx api.Context, params *sodor.AlertPluginInstance) (*sodor.AlertPluginReply, *api.JsonRpcError) {
	if params.Id == 0 {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("invalid params of alert_plugin_instance"))
	}

	exist, err := metastore.GetInstance().AlertPluginInstanceExist(params)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if !exist {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("alert_plugin_instance not exist"))
	}

	use, err := metastore.GetInstance().AlertPluginInstanceUsedInAlertGroup(params)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if use {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("alert_plugin_instance is still using in alert_group"))
	}

	if err = metastore.GetInstance().DeleteAlertPluginInstance(params); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	ctx.Log("DeleteAlertPluginInstance Done: %+v", params)
	return &sodor.AlertPluginReply{Id: params.Id}, nil
}
