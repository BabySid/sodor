package thomas

import (
	"context"
	"errors"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"sodor/fat_controller/config"
	"strconv"
)

type Thomas struct {
	Host string
	Port int
}

func (t *Thomas) HandShake(id int32) error {
	conn, err := t.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := sodor.NewThomasClient(conn)

	var req sodor.FatCtrlInfo
	req.Id = id
	req.Name = config.GetInstance().AppName
	req.Version = config.GetInstance().AppVersion
	req.Host = config.GetInstance().LocalIP
	req.Port = int32(config.GetInstance().Port)

	_, err = cli.HandShake(context.Background(), &req)
	if s, ok := status.FromError(err); ok {
		if s != nil {
			log.Warnf("HandShake to thomas failed. code = %d, msg = %s", s.Code(), s.Message())
			return errors.New(s.String())
		}
	}

	return err
}

func (t *Thomas) RunTask(jobIns int32, taskIns int32, task *sodor.Task) error {
	conn, err := t.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := sodor.NewThomasClient(conn)
	var req sodor.RunTaskRequest
	req.TaskInstanceId = taskIns
	req.JobInstanceId = jobIns
	req.JobId = task.JobId
	req.TaskId = task.Id
	req.Task = task

	_, err = cli.RunTask(context.Background(), &req)
	if s, ok := status.FromError(err); ok {
		if s != nil {
			log.Warnf("RunTask to thomas failed. code = %d, msg = %s", s.Code(), s.Message())
			return errors.New(s.String())
		}
	}

	return err
}

func (t *Thomas) dial() (*grpc.ClientConn, error) {
	host := t.Host + ":" + strconv.Itoa(t.Port)
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Warnf("Dial host = %s failed. err = %s", host, err)
		return nil, err
	}

	return conn, nil
}
