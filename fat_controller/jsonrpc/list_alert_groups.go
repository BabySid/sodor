package jsonrpc

import (
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ListAlertGroup(ctx api.Context, _ *interface{}) (*sodor.AlertGroups, *api.JsonRpcError) {
	ags, err := metastore.GetInstance().ListAlertGroups()
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}
	ctx.Log("ListAlertGroup Done: %d", len(ags.GetAlertGroups()))
	return ags, nil
}
