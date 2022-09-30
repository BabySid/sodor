package scheduler

import (
	log "github.com/sirupsen/logrus"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/thomas"
)

func handShakeWithOverDueThomas() {
	ts, err := metastore.GetInstance().SelectInvalidThomas()
	if err != nil {
		log.Warnf("SelectInvalidThomas return err=%s", err)
		return
	}

	for _, t := range ts {
		pingThomas(&t)
	}
}

func pingThomas(thomas *metastore.Thomas) {
	err := PingThomas(int32(thomas.ID), thomas.Host, thomas.Port)
	if err != nil {
		log.Warnf("PingThomas failed. err=%s", err)
	}
}

func PingThomas(id int32, host string, port int) error {
	t := thomas.Thomas{
		Host: host,
		Port: port,
	}

	err := t.HandShake()

	if err != nil {
		return metastore.GetInstance().UpdateThomasStatus(id, "")
	}
	return nil
}
