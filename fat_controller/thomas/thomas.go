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

func (t *Thomas) HandShake() error {
	host := t.Host + ":" + strconv.Itoa(t.Port)
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Warnf("Dial host = %s failed. err = %s", host, err)
		return err
	}
	defer conn.Close()

	cli := sodor.NewThomasClient(conn)

	var req sodor.FatCtrlInfo
	// todo build req from metastore
	req.Id = 1
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
