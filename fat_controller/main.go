package main

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc"
	"github.com/BabySid/gorpc/http/httpcfg"
	logOption "github.com/BabySid/gorpc/log"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"path/filepath"
	"sodor/base"
	"sodor/fat_controller/config"
	"sodor/fat_controller/grpc"
	"sodor/fat_controller/jsonrpc"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
	"sort"
	"syscall"
)

var (
	AppVersion string
	AppName    = filepath.Base(os.Args[0])
	server     *gorpc.Server
)

func main() {
	ss := gobase.NewSignalSet()
	ss.Register(syscall.SIGTERM, exit)

	app := NewApp()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.UseShortOptionHandling = true
	app.Name = AppName
	app.Version = AppVersion
	app.Usage = "fat_ctrl is responsible for ensuring that all trains arrive at the station on time and play their due role"

	app.Action = runApp

	app.Flags = config.GlobalFlags
	app.Commands = []*cli.Command{}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	cli.AppHelpTemplate = base.AppHelpTemplate
	originalHelpPrinter := cli.HelpPrinter
	cli.HelpPrinter = func(w io.Writer, tmpl string, data interface{}) {
		originalHelpPrinter(w, tmpl, base.HelpData{
			App:        data,
			FlagGroups: config.AppHelpFlagGroups,
		})
	}

	return app
}

func runApp(ctx *cli.Context) error {
	if ctx.NumFlags() == 0 {
		cli.ShowAppHelpAndExit(ctx, 1)
	}

	err := initComponent(ctx)
	if err != nil {
		return err
	}

	if ctx.Bool(config.InitMetaStore.Name) {
		initializeMetaStore()
		return nil
	}

	server = gorpc.NewServer(httpcfg.ServerOption{Codec: httpcfg.ProtobufCodec})
	_ = server.RegisterJsonRPC("rpc", &jsonrpc.Service{})
	_ = server.RegisterGrpc(&sodor.FatController_ServiceDesc, &grpc.Service{})

	var rotator *logOption.Rotator
	if !ctx.Bool(config.DebugMode.Name) {
		rotator = &logOption.Rotator{
			LogMaxAge: ctx.Int(config.LogMaxAge.Name),
			LogPath:   ctx.String(config.LogPath.Name),
		}
	}

	return server.Run(gorpc.ServerOption{
		Addr:        ctx.String(config.ListenAddr.Name),
		ClusterName: "fat_ctrl",
		Rotator:     rotator,
		LogLevel:    ctx.String(config.LogLevel.Name),
	})
}

func exit(sig os.Signal) {
	log.Infof("%s exit by recving the signal %v", AppName, sig)
	_ = server.Stop()
	os.Exit(0)
}

func initializeMetaStore() {
	err := metastore.GetInstance().AutoMigrate()
	if err != nil {
		log.Warnf("init meta automigrate failed. err=%s", err)
		return
	}
}

func initComponent(ctx *cli.Context) error {
	err := config.GetInstance().InitFromFlags(ctx)
	if err != nil {
		log.Fatalf("config init failed. err=%s", err)
	}
	config.GetInstance().AppName = AppName
	config.GetInstance().AppVersion = AppVersion

	_ = metastore.GetInstance()

	err = scheduler.GetInstance().Start()
	if err != nil {
		log.Fatalf("scheduler init failed. err=%s", err)
	}
	return nil
}
