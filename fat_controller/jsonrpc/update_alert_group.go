package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
)

func (s *Service) UpdateAlertGroup(ctx *httpapi.APIContext, params *sodor.AlertGroup) (*sodor.AlertGroupReply, *httpapi.JsonRpcError) {
	return nil, nil
}
