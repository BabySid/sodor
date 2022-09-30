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
	"sodor/thomas/grpc"
	"sodor/thomas/task"
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

	app.Flags = globalFlags
	app.Commands = []*cli.Command{
		task.ShellCommand,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	cli.AppHelpTemplate = base.AppHelpTemplate
	originalHelpPrinter := cli.HelpPrinter
	cli.HelpPrinter = func(w io.Writer, tmpl string, data interface{}) {
		originalHelpPrinter(w, tmpl, base.HelpData{
			App:        data,
			FlagGroups: appHelpFlagGroups,
		})
	}

	return app
}

func runApp(ctx *cli.Context) error {
	if ctx.NumFlags() == 0 {
		cli.ShowAppHelpAndExit(ctx, 1)
	}

	if ctx.Bool(standalone.Name) {
		runStandalone(ctx)
		return nil
	}

	var rotator *logOption.Rotator
	if !ctx.Bool(debugMode.Name) {
		rotator = &logOption.Rotator{
			LogMaxAge: ctx.Int(logMaxAge.Name),
			LogPath:   ctx.String(logPath.Name),
		}
	}

	server = gorpc.NewServer(httpcfg.DefaultOption)
	_ = server.RegisterGrpc(&sodor.Thomas_ServiceDesc, grpc.NewService())

	return server.Run(gorpc.ServerOption{
		Addr:        ctx.String(listenAddr.Name),
		ClusterName: "thomas",
		Rotator:     rotator,
		LogLevel:    ctx.String(logLevel.Name),
	})
}

func exit(sig os.Signal) {
	log.Infof("%s exit by recving the signal %v", AppName, sig)
	_ = server.Stop()
	os.Exit(0)
}

func runStandalone(ctx *cli.Context) {

}
