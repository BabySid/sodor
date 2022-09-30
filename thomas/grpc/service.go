package grpc

import (
	"errors"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"sync"
)

type Service struct {
	sodor.UnimplementedThomasServer
	thomasID int32

	mux          sync.Mutex
	fatCtrlHosts []host
	fatCtrlIdx   int
}

func NewService() *Service {
	return &Service{
		thomasID:     0,
		fatCtrlHosts: make([]host, 0),
		fatCtrlIdx:   0,
	}
}

type host struct {
	IP   string
	port int
}

func (s *Service) getFatCtrlHost() (*grpc.ClientConn, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if len(s.fatCtrlHosts) == 0 {
		return nil, errors.New("don't known fat_ctrl's address")
	}

	initIdx := s.fatCtrlIdx
	for {
		h := s.fatCtrlHosts[s.fatCtrlIdx].IP + ":" + strconv.Itoa(s.fatCtrlHosts[s.fatCtrlIdx].port)
		conn, err := grpc.Dial(h, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Warnf("Dial host = %s failed. err = %s", h, err)
			continue
		}
		s.fatCtrlIdx = (s.fatCtrlIdx + 1) % len(s.fatCtrlHosts)
		if s.fatCtrlIdx == initIdx {
			break
		}
		return conn, nil
	}

	return nil, errors.New("cannot found valid fat_ctrl's address")
}

func (s *Service) updateFatCtrlHost(ip string, port int) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, h := range s.fatCtrlHosts {
		if h.IP == ip && h.port == port {
			return nil
		}
	}

	s.fatCtrlHosts = append(s.fatCtrlHosts, host{
		IP:   ip,
		port: port,
	})

	return nil
}
