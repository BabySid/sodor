package jsonrpc

import (
	"context"
	"github.com/BabySid/gorpc/http/httpapi"
	"sodor/fat_controller/grpc"
)
import "github.com/BabySid/proto/sodor"

func (s *Service) CreateJob(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.Reply, error) {
	gs := &grpc.Service{}
	ctx.ToLog("CreateJob: %+v", params)
	reply, err := gs.CreateJob(context.Background(), params)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
