package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
)

func (s *Service) ShowAlertPlugins(ctx *httpapi.APIContext, _ *interface{}) (*sodor.AlertPlugin, *httpapi.JsonRpcError) {
	return nil, nil
}
