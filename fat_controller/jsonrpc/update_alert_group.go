package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) UpdateAlertGroup(ctx *httpapi.APIContext, params *sodor.AlertGroup) (*sodor.AlertGroupReply, *httpapi.JsonRpcError) {
	if err := checkAlertGroupValid(params, false); err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, err)
	}

	exist, err := metastore.GetInstance().AlertGroupExist(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	if !exist {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("alert_group not exist"))
	}

	for _, ins := range params.PluginInstance {
		exist, err = metastore.GetInstance().AlertPluginInstanceExist(&sodor.AlertPluginInstance{Id: ins})
		if err != nil {
			return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
		}

		if exist {
			return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("plugin_instance not exist"))
		}
	}

	if err = metastore.GetInstance().UpdateAlertGroup(params); err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	ctx.ToLog("UpdateAlertGroup Done: %+v", params)
	return &sodor.AlertGroupReply{Id: params.Id}, nil
}
