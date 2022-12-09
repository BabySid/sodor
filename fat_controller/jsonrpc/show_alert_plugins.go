package jsonrpc

import (
	"github.com/BabySid/gorpc/http/httpapi"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/alert"
)

func (s *Service) ShowAlertPlugins(ctx *httpapi.APIContext, _ *interface{}) (*sodor.AlertPlugins, *httpapi.JsonRpcError) {
	alerts := alert.GetInstance().GetAlertPlugins()

	var aps sodor.AlertPlugins
	aps.AlertPlugins = make([]*sodor.AlertPlugin, len(alerts))

	for i, a := range alerts {
		aps.AlertPlugins[i] = &sodor.AlertPlugin{
			Name:   a.GetName(),
			Params: a.GetParams(),
		}
	}

	ctx.ToLog("ShowAlertPlugins Done: %d", len(alerts))
	return &aps, nil
}
