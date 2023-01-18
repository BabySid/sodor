package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/alert"
	"sodor/fat_controller/metastore"
)

func (s *Service) CreateAlertGroup(ctx *httpapi.APIContext, params *sodor.AlertGroup) (*sodor.AlertGroupReply, *httpapi.JsonRpcError) {
	if err := checkAlertGroupValid(params, true); err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, err)
	}

	exist, err := metastore.GetInstance().AlertGroupExist(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	if exist {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("alert_group exist"))
	}

	for _, ins := range params.PluginInstances {
		exist, err = metastore.GetInstance().AlertPluginInstanceExist(&sodor.AlertPluginInstance{Id: ins})
		if err != nil {
			return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
		}

		if !exist {
			return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("plugin_instance not exist"))
		}
	}

	if err = metastore.GetInstance().InsertAlertGroup(params); err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	if params.Name == alert.SystemAlertGroupName {
		if err = alert.GetInstance().ResetAlertGroupID(); err != nil {
			return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
		}
	}
	ctx.ToLog("CreateAlertGroup Done: %+v", params)
	return &sodor.AlertGroupReply{Id: params.Id}, nil
}
