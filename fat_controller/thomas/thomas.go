package thomas

import (
	"context"
	"errors"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"sodor/fat_controller/util"
	"strconv"
)

type Thomas struct {
	Host string
	Port int
}

func (t *Thomas) HandShake(id int32) (*sodor.ThomasInfo, error) {
	conn, err := t.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := sodor.NewThomasClient(conn)

	var req sodor.HandShakeWithThomasRequest
	req.Thomas = &sodor.ThomasInfo{Id: id}
	req.FatCtrls = util.BuildFatCtrlInfos()

	resp, err := cli.HandShake(context.Background(), &req)
	if s, ok := status.FromError(err); ok {
		if s != nil {
			log.Warnf("HandShake to thomas failed. code=%d, msg=%s", s.Code(), s.Message())
			return nil, errors.New(s.String())
		}
	}

	return resp, err
}

func (t *Thomas) RunTask(task *sodor.Task, ins *sodor.TaskInstance) error {
	conn, err := t.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := sodor.NewThomasClient(conn)
	var req sodor.RunTaskRequest
	req.Task = task
	req.TaskInstance = ins

	_, err = cli.RunTask(context.Background(), &req)
	if s, ok := status.FromError(err); ok {
		if s != nil {
			log.Warnf("RunTask to thomas failed. code=%d, msg=%s", s.Code(), s.Message())
			return errors.New(s.String())
		}
	}

	return err
}

func (t *Thomas) dial() (*grpc.ClientConn, error) {
	host := t.Host + ":" + strconv.Itoa(t.Port)
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Warnf("Dial host=%s failed. err=%s", host, err)
		return nil, err
	}

	return conn, nil
}
