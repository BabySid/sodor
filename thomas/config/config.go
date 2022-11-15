package config

import (
	"github.com/BabySid/gobase"
	"github.com/urfave/cli/v2"
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

	port, err := strconv.Atoi(arr[0])
	gobase.TrueF(err == nil, "invalid port of %s", ListenAddr.Name)
	c.Port = port

	c.RetryInterval = time.Second * 5
	return nil
}
