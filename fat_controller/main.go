package main

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc"
	logOption "github.com/BabySid/gorpc/log"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"path/filepath"
	"sodor/base"
	"sodor/fat_controller/jsonrpc"
	"sort"
	"syscall"
)

var (
	AppVersion string
	AppName    = filepath.Base(os.Args[0])
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

	app.Flags = []cli.Flag{
		listenAddr,
		metaStore,
		logLevel,
		logPath,
		logMaxAge,
		debugMode,
	}
	app.Commands = []*cli.Command{}

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

	s := gorpc.NewServer()
	s.RegisterJsonRPC("rpc", &jsonrpc.Service{})

	var rotator *logOption.Rotator
	if !ctx.Bool(debugMode.Name) {
		rotator = &logOption.Rotator{
			LogMaxAge: ctx.Int(logMaxAge.Name),
			LogPath:   ctx.String(logPath.Name),
		}
	}

	return s.Run(gorpc.ServerOption{
		Addr:        ctx.String(listenAddr.Name),
		ClusterName: "fat_ctrl",
		Rotator:     rotator,
		LogLevel:    ctx.String(logLevel.Name),
	})
}

func exit(sig os.Signal) {
	log.Infof("%s exit by recving the signal %v", AppName, sig)
	os.Exit(0)
}
