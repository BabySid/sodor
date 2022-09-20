package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
)

func (s *Service) CreateJob(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.JobReply, error) {
	ctx.ToLog("CreateJob: %+v", params)
	return &sodor.JobReply{Id: 1}, nil
}
