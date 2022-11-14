package fat_ctrl

import (
	"context"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

func (fc *FatCtrl) UpdateTaskInstance(ins *sodor.TaskInstance) error {
	conn, err := fc.getFatCtrlConn()
	if err != nil {
		log.Warnf("getFatCtrlConn failed. err=%s", err)
		return err
	}

	defer conn.Close()

	cli := sodor.NewFatControllerClient(conn)
	_, err = cli.UpdateTaskInstance(context.Background(), ins)

	if s, ok := status.FromError(err); ok {
		if s != nil {
			log.Warnf("UpdateTaskInstance to fat_ctrl failed. code=%d, msg=%s", s.Code(), s.Message())
		}
	} else {
		if err != nil {
			log.Warnf("UpdateTaskInstance to fat_ctrl failed. err=%s", err)
		}
	}

	return err
}
