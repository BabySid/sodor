package grpc

import (
	"github.com/BabySid/proto/sodor"
)

type Service struct {
	sodor.UnimplementedFatControllerServer
}
