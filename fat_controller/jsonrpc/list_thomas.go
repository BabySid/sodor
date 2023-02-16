package jsonrpc

import (
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
)

func (s *Service) ListThomas(ctx api.Context, _ *interface{}) (*sodor.ThomasInfos, *api.JsonRpcError) {
	thomas, err := metastore.GetInstance().ListAllThomas()
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}
	ctx.Log("ListThomas Done: %d", len(thomas.GetThomasInfos()))
	return thomas, nil
}
