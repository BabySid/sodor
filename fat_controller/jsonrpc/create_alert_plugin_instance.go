package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) CreateAlertPluginInstance(ctx *httpapi.APIContext, params *sodor.AlertPluginInstance) (*sodor.AlertPluginReply, *httpapi.JsonRpcError) {
	if err := checkAlertPluginValid(params, true); err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, err)
	}

	exist, err := metastore.GetInstance().AlertPluginInstanceExist(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	if exist {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("alert_plugin exist"))
	}

	if err = metastore.GetInstance().InsertAlertPluginInstance(params); err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	ctx.ToLog("CreateAlertPluginInstance Done: %+v", params)
	return &sodor.AlertPluginReply{Id: params.Id}, nil
}
