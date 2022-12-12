package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ListAlertPluginInstances(ctx *httpapi.APIContext, _ *interface{}) (*sodor.AlertPluginInstances, *httpapi.JsonRpcError) {
	aps, err := metastore.GetInstance().ListAlertPluginInstances()

	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	ctx.ToLog("ListAlertPluginInstances Done: %d", len(aps.AlertPluginInstances))
	return aps, nil
}
