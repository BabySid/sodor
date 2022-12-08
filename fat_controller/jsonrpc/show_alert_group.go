package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
)

func (s *Service) ShowAlertGroup(ctx *httpapi.APIContext, params *sodor.AlertGroup) (*sodor.AlertGroup, *httpapi.JsonRpcError) {
	return nil, nil
}
