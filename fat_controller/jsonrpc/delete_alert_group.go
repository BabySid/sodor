package jsonrpc

import (
	"errors"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/alert"
	"sodor/fat_controller/metastore"
)

func (s *Service) DeleteAlertGroup(ctx api.Context, params *sodor.AlertGroup) (*sodor.AlertGroupReply, *api.JsonRpcError) {
	if params.Id == 0 {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("invalid params of alert_group"))
	}

	exist, err := metastore.GetInstance().AlertGroupExist(params)
	if err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if !exist {
		return nil, api.NewJsonRpcErrFromCode(api.InvalidParams, errors.New("alert_group not exist"))
	}

	if err = metastore.GetInstance().DeleteAlertGroup(params); err != nil {
		return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
	}

	if params.Id == alert.GetInstance().AlertGroupID() {
		if err = alert.GetInstance().ResetAlertGroupID(); err != nil {
			return nil, api.NewJsonRpcErrFromCode(api.InternalError, err)
		}
	}

	ctx.Log("DeleteAlertGroup Done: %+v", params)
	return &sodor.AlertGroupReply{Id: params.Id}, nil
}
