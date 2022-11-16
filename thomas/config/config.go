package config

import (
	"github.com/BabySid/gobase"
	"github.com/urfave/cli/v2"
	"sodor/base"
	"strconv"
	"strings"
	"sync"
	"time"
)

type config struct {
	LocalIP string
	Port    int

	DataPath      string
	TaskIdentity  string
	AppName       string
	AppVersion    string
	RetryInterval time.Duration
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

	port, err := strconv.Atoi(arr[1])
	gobase.TrueF(err == nil, "invalid port of %s", addr)
	c.Port = port

	c.LocalIP = base.LocalHost

	c.RetryInterval = time.Second * 5
	return nil
}
