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
	if err := checkJobValid(params, false); err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, err)
	}

	exist, err := metastore.GetInstance().JobExist(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	if !exist {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("job not exist"))
	}

	if params.AlertGroupId > 0 {
		exist, err = metastore.GetInstance().AlertGroupExist(&sodor.AlertGroup{Id: params.AlertGroupId})
		if err != nil {
			return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
		}

		if !exist {
			return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("alert_group not exist"))
		}
	}

	err = metastore.GetInstance().UpdateJob(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	err = scheduler.GetInstance().Remove(params)
	gobase.True(err == nil)

	if params.ScheduleMode == sodor.ScheduleMode_ScheduleMode_Crontab {
		err = scheduler.GetInstance().AddJob(params)
		gobase.True(err == nil)
	}

	ctx.ToLog("UpdateJob Done: %d", params.Id)
	return &sodor.JobReply{Id: params.Id}, nil
}
