package main

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/gorpc/codec"
	"github.com/BabySid/proto/sodor"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"os"
	"path/filepath"
	"sodor/thomas/config"
	"sodor/thomas/grpc"
	"sodor/thomas/routine"
	"sodor/thomas/task_runner"
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
	app.Name = filepath.Base(os.Args[0])
	app.Version = AppVersion
	app.Usage = "thomas: a famous little tank engine run job from fat_controller"

	app.Action = runApp

	app.Flags = config.GlobalFlags
	app.Before = func(ctx *cli.Context) error {
		BeforeFunc := altsrc.InitInputSourceWithContext(config.GlobalFlags, altsrc.NewTomlSourceFromFlagFunc(config.ConfFile.Name))
		_ = BeforeFunc(ctx)
		return nil
	}
	app.Commands = []*cli.Command{
		&runTask,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	//cli.AppHelpTemplate = base.AppHelpTemplate
	//originalHelpPrinter := cli.HelpPrinter
	//cli.HelpPrinter = func(w io.Writer, tmpl string, data interface{}) {
	//	originalHelpPrinter(w, tmpl, base.HelpData{
	//		App:        data,
	//		FlagGroups: config.AppHelpFlagGroups,
	//	})
	//}

	return app
}

func runApp(ctx *cli.Context) error {
	var rotator *api.Rotator
	if !ctx.Bool(config.DebugMode.Name) {
		rotator = &api.Rotator{
			LogMaxAge: ctx.Int(config.LogMaxAge.Name),
			LogPath:   ctx.String(config.LogPath.Name),
		}
	}

	server = gorpc.NewServer(api.ServerOption{
		Addr:        ctx.String(config.ListenAddr.Name),
		ClusterName: "thomas",
		Rotator:     rotator,
		LogLevel:    ctx.String(config.LogLevel.Name),
		Codec:       codec.JsonCodec,
		BeforeRun: func() error {
			return initComponent(ctx)
		},
	})

	_ = server.RegisterGrpc(&sodor.Thomas_ServiceDesc, &grpc.Service{})

	return server.Run()
}

func initComponent(ctx *cli.Context) error {
	if err := config.GetInstance().InitFromFlags(ctx); err != nil {
		log.Fatalf("config init failed. err=%s", err)
	}

	config.GetInstance().AppName = AppName
	config.GetInstance().AppVersion = AppVersion

	if err := routine.GetInstance().Start(); err != nil {
		log.Fatalf("routine start failed. err=%s", err)
	}

	if err := task_runner.GetTaskEnv().LoadTasksStatus(); err != nil {
		log.Fatalf("task_runner load status failed. err=%s", err)
	}

	return nil
}

func exit(sig os.Signal) {
	log.Infof("%s exit by recving the signal %v", AppName, sig)
	_ = server.Stop()
	os.Exit(0)
}
