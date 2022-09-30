package jsonrpc

import (
	"errors"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
)

func (s *Service) UpdateJob(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.JobReply, *httpapi.JsonRpcError) {
	if err := checkTaskValid(params, false); err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams, httpapi.SysCodeMap[httpapi.InvalidParams], err)
	}

	exist, err := metastore.GetInstance().JobExist(params)
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	if !exist {
		return nil, httpapi.NewJsonRpcError(httpapi.InvalidParams,
			httpapi.SysCodeMap[httpapi.InvalidParams], errors.New("job not exist"))
	}

	err = metastore.GetInstance().UpdateJob(params)
	if err != nil {
		return nil, httpapi.NewJsonRpcError(httpapi.InternalError, httpapi.SysCodeMap[httpapi.InternalError], err)
	}

	err = scheduler.GetInstance().Remove(params)
	gobase.True(err == nil)

	if params.ScheduleMode == sodor.ScheduleMode_ScheduleMode_Crontab {
		err = scheduler.GetInstance().AddJob(params)
		gobase.True(err == nil)
	}

	ctx.ToLog("UpdateJob Done: %+v", params)
	return &sodor.JobReply{Id: params.Id}, nil
}
