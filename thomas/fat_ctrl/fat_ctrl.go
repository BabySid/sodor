package fat_ctrl

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"sync"
	"time"
)

type FatCtrl struct {
	thomasID int32

	mux          sync.Mutex
	fatCtrlHosts []host
	fatCtrlIdx   int

	startTime int32
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
			thomasID:     0,
			fatCtrlHosts: make([]host, 0),
			fatCtrlIdx:   0,
			startTime:    int32(time.Now().Unix()),
		}
	})
	return singleton
}

func (fc *FatCtrl) Run() {
	fc.HandShake()
}

func (fc *FatCtrl) getFatCtrlConn() (*grpc.ClientConn, error) {
	fc.mux.Lock()
	defer fc.mux.Unlock()

	if fc.thomasID == 0 {
		return nil, errors.New("thomas_id is not set")
	}

	if len(fc.fatCtrlHosts) == 0 {
		return nil, errors.New("don't known fat_ctrl's address")
	}

	initIdx := fc.fatCtrlIdx
	for {
		h := fc.fatCtrlHosts[fc.fatCtrlIdx].IP + ":" + strconv.Itoa(fc.fatCtrlHosts[fc.fatCtrlIdx].port)
		conn, err := grpc.Dial(h, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (fc *FatCtrl) UpdateFatCtrlHost(ip string, port int) error {
	fc.mux.Lock()
	defer fc.mux.Unlock()

	for _, h := range fc.fatCtrlHosts {
		if h.IP == ip && h.port == port {
			return nil
		}
	}

	fc.fatCtrlHosts = append(fc.fatCtrlHosts, host{
		IP:   ip,
		port: port,
	})

	log.Infof("add new fat_ctrl's address. host=%s port=%d", ip, port)
	return nil
}

func (fc *FatCtrl) SetThomasID(id int32) {
	fc.thomasID = id
}
