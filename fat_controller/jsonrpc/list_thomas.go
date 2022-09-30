package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ListThomas(ctx *httpapi.APIContext, params *interface{}) (*sodor.ThomasInfos, *httpapi.JsonRpcError) {
	thomas, err := metastore.GetInstance().ListAllThomas()
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}
	ctx.ToLog("ListThomas Done: %d", len(thomas.GetThomasInfos()))
	return thomas, nil
}
