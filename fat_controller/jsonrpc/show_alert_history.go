package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
)

func (s *Service) ShowAlertGroupHistory(ctx *httpapi.APIContext, params *sodor.AlertHistory) (*sodor.AlertHistory, *httpapi.JsonRpcError) {
	return nil, nil
}
