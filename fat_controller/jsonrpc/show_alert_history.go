package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ShowAlertPluginInstanceHistories(
	ctx *httpapi.APIContext,
	params *sodor.AlertPluginInstanceHistory) (*sodor.AlertPluginInstanceHistories, *httpapi.JsonRpcError) {
	if params.GroupId == 0 && params.InstanceId == 0 {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("invalid params of alert_plugin_instance_history"))
	}

	rs, err := metastore.GetInstance().ShowAlertPluginInstanceHistories(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	ctx.ToLog("ShowAlertPluginInstanceHistories Done: %d", len(rs.AlertPluginInstanceHistory))
	return rs, nil
}
