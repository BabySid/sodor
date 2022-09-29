package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ListThomas(ctx *httpapi.APIContext, params *interface{}) (*sodor.ThomasInstances, *httpapi.JsonRpcError) {
	thomas, err := metastore.GetInstance().ListAllThomas()
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}
	ctx.ToLog("ListJobs Done: %d", len(thomas.GetThomasInstances()))
	return thomas, nil
}
