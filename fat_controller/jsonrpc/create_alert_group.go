package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/alert"
	"sodor/fat_controller/metastore"
)

func (s *Service) CreateAlertGroup(ctx api.Context, params *sodor.AlertGroup) (*sodor.AlertGroupReply, *api.JsonRpcError) {
	if err := checkAlertGroupValid(params, true); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, err)
	}

	exist, err := metastore.GetInstance().AlertGroupExist(params)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if exist {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("alert_group exist"))
	}

	for _, ins := range params.PluginInstances {
		exist, err = metastore.GetInstance().AlertPluginInstanceExist(&sodor.AlertPluginInstance{Id: ins})
		if err != nil {
			return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
		}

		if !exist {
			return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("plugin_instance not exist"))
		}
	}

	if err = metastore.GetInstance().InsertAlertGroup(params); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if params.Name == alert.SystemAlertGroupName {
		if err = alert.GetInstance().ResetAlertGroupID(); err != nil {
			return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
		}
	}
	ctx.Log("CreateAlertGroup Done: %+v", params)
	return &sodor.AlertGroupReply{Id: params.Id}, nil
}
