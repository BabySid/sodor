package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) DeleteAlertPluginInstance(ctx *httpapi.APIContext, params *sodor.AlertPluginInstance) (*sodor.AlertPluginReply, *httpapi.JsonRpcError) {
	if params.Id == 0 {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("invalid params of alert_plugin_instance"))
	}

	exist, err := metastore.GetInstance().AlertPluginInstanceExist(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	if !exist {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("alert_plugin_instance not exist"))
	}

	use, err := metastore.GetInstance().AlertPluginInstanceUsedInAlertGroup(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	if use {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("alert_plugin_instance is still using in alert_group"))
	}

	if err = metastore.GetInstance().DeleteAlertPluginInstance(params); err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	ctx.ToLog("DeleteAlertPluginInstance Done: %+v", params)
	return &sodor.AlertPluginReply{Id: params.Id}, nil
}
