package config

import (
	"github.com/BabySid/gobase"
	"github.com/urfave/cli/v2"
	"sodor/base"
	"strconv"
	"strings"
	"sync"
)

type config struct {
	LocalIP           string
	Port              int
	MetaStoreUri      string
	MaxThomasInstance uint32
	MaxJobInstance    uint32
	MaxAlertHistory   uint32

	AppName    string
	AppVersion string
}

var (
	once      sync.Once
	singleton *config
)

func GetInstance() *config {
	once.Do(func() {
		singleton = &config{}
	})
	return singleton
}

func (c *config) InitFromFlags(ctx *cli.Context) error {
	addr := ctx.String(ListenAddr.Name)

	arr := strings.Split(addr, ":")
	gobase.TrueF(len(arr) >= 1, "%s format is '[$host]:$port'", ListenAddr.Name)

	c.LocalIP = base.LocalHost
	port, err := strconv.Atoi(arr[1])
	gobase.TrueF(err == nil, "invalid port of %s", addr)
	c.Port = port

	c.MaxThomasInstance = 64
	c.MaxJobInstance = 64
	c.MaxAlertHistory = 64

	return nil
}
