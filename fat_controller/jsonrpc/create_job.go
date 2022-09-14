package jsonrpc

import "github.com/BabySid/gorpc/http/httpapi"

func (s *Service) CreateJob(ctx *httpapi.APIContext, params *interface{}) (interface{}, error) {
	ctx.ToLog("CreateJob")
	return "OK", nil
}
