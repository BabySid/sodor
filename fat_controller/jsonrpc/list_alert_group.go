package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ListAlertGroup(ctx *httpapi.APIContext, _ *interface{}) (*sodor.AlertGroups, *httpapi.JsonRpcError) {
	ags, err := metastore.GetInstance().ListAlertGroups()
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}
	ctx.ToLog("ListAlertGroup Done: %d", len(ags.GetAlertGroups()))
	return ags, nil
}
