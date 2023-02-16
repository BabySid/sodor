package jsonrpc

import (
	"errors"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
)

func (s *Service) CreateJob(ctx api.Context, params *sodor.Job) (*sodor.JobReply, *api.JsonRpcError) {
	if err := checkJobValid(params, true); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, err)
	}

	exist, err := metastore.GetInstance().JobExist(params)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if exist {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("job exist"))
	}

	if params.AlertGroupId > 0 {
		exist, err = metastore.GetInstance().AlertGroupExist(&sodor.AlertGroup{Id: params.AlertGroupId})
		if err != nil {
			return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
		}

		if !exist {
			return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("alert_group not exist"))
		}
	}

	err = metastore.GetInstance().InsertJob(params)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if params.AlertGroupId > 0 {
		exist, err = metastore.GetInstance().AlertGroupExist(&sodor.AlertGroup{Id: params.AlertGroupId})
		if err != nil {
			return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
		}

		if !exist {
			return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("alert_group not exist"))
		}
	}

	if params.ScheduleMode == sodor.ScheduleMode_ScheduleMode_Crontab {
		err = scheduler.GetInstance().AddJob(params)
		gobase.True(err == nil)
	}

	ctx.Log("CreateJob Done: %+v", params)
	return &sodor.JobReply{Id: params.Id}, nil
}
