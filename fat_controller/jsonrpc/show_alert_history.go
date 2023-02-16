package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ShowAlertPluginInstanceHistories(
	ctx api.Context,
	params *sodor.AlertPluginInstanceHistory) (*sodor.AlertPluginInstanceHistories, *api.JsonRpcError) {
	if params.GroupId == 0 && params.InstanceId == 0 {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("invalid params of alert_plugin_instance_history"))
	}

	rs, err := metastore.GetInstance().ShowAlertPluginInstanceHistories(params)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	ctx.Log("ShowAlertPluginInstanceHistories Done: %d", len(rs.AlertPluginInstanceHistory))
	return rs, nil
}
