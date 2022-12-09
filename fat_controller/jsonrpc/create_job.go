package jsonrpc

import (
	"errors"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
)

func (s *Service) CreateJob(ctx *httpapi.APIContext, params *sodor.Job) (*sodor.JobReply, *httpapi.JsonRpcError) {
	if err := checkJobValid(params, true); err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, err)
	}

	exist, err := metastore.GetInstance().JobExist(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
	}

	if exist {
		return nil, httpapi.NewJRpcErr(httpapi.InvalidParams, errors.New("job exist"))
	}

	err = metastore.GetInstance().InsertJob(params)
	if err != nil {
		return nil, httpapi.NewJRpcErr(httpapi.InternalError, err)
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

	if params.ScheduleMode == sodor.ScheduleMode_ScheduleMode_Crontab {
		err = scheduler.GetInstance().AddJob(params)
		gobase.True(err == nil)
	}

	ctx.ToLog("CreateJob Done: %+v", params)
	return &sodor.JobReply{Id: params.Id}, nil
}
