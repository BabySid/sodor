package fat_ctrl

import (
	"errors"
	"github.com/BabySid/gorpc"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"sodor/thomas/config"
	"strconv"
	"sync"
)

type FatCtrl struct {
	mux          sync.Mutex
	fatCtrlHosts []host
	fatCtrlIdx   int
}

type host struct {
	IP   string
	port int
}

var (
	once      sync.Once
	singleton *FatCtrl
)

func GetInstance() *FatCtrl {
	once.Do(func() {
		singleton = &FatCtrl{
			fatCtrlHosts: make([]host, 0),
			fatCtrlIdx:   0,
		}
	})
	return singleton
}

func (fc *FatCtrl) Run() {
	fc.HandShake()
}

func (fc *FatCtrl) getFatCtrlConn() (api.Client, error) {
	fc.mux.Lock()
	defer fc.mux.Unlock()

	if config.GetInstance().ThomasID == 0 {
		return nil, errors.New("thomas_id is not set")
	}

	if len(fc.fatCtrlHosts) == 0 {
		return nil, errors.New("don't known fat_ctrl's address")
	}

	initIdx := fc.fatCtrlIdx
	for {
		h := "grpc://" + fc.fatCtrlHosts[fc.fatCtrlIdx].IP + ":" + strconv.Itoa(fc.fatCtrlHosts[fc.fatCtrlIdx].port)
		conn, err := gorpc.Dial(h, api.ClientOption{})
		fc.fatCtrlIdx = (fc.fatCtrlIdx + 1) % len(fc.fatCtrlHosts)
		if err == nil {
			return conn, nil
		}
		log.Warnf("Dial thomas(host=%s) failed. err=%s", h, err)
		if fc.fatCtrlIdx == initIdx {
			break
		}
	}

	return nil, errors.New("cannot found valid fat_ctrl's address")
}

func (fc *FatCtrl) UpdateFatCtrlHost(infos *sodor.FatCtrlInfos) error {
	fc.mux.Lock()
	defer fc.mux.Unlock()

	fc.fatCtrlHosts = make([]host, len(infos.FatCtrlInfos))
	for i, info := range infos.FatCtrlInfos {
		fc.fatCtrlHosts[i] = host{
			IP:   info.Host,
			port: int(info.Port),
		}
	}
	return nil
}
