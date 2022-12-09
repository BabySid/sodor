package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ShowAlertGroupInstance(ctx *httpapi.APIContext, params *sodor.AlertGroup) (*sodor.AlertGroupInstances, *httpapi.JsonRpcError) {
	rs, err := metastore.GetInstance().ShowAlertGroupInstanceByGroupID(params.Id)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}
	ctx.ToLog("ShowAlertGroupInstance Done: %d", len(rs.AlertGroupInstances))
	return rs, nil
}
