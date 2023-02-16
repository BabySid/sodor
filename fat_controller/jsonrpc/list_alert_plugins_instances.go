package jsonrpc

import (
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ListAlertPluginInstances(ctx api.Context, _ *interface{}) (*sodor.AlertPluginInstances, *api.JsonRpcError) {
	aps, err := metastore.GetInstance().ListAlertPluginInstances()

	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	ctx.Log("ListAlertPluginInstances Done: %d", len(aps.AlertPluginInstances))
	return aps, nil
}
