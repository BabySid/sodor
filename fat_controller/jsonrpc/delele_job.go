package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
)

func (s *Service) DeleteJob(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.JobReply, *httpapi.JsonRpcError) {
	return nil, nil
}
