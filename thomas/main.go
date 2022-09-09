package main

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"path/filepath"
	"sodor/base"
	"sodor/thomas/task"
	"sort"
	"syscall"
)

// todo add signal handler
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
	app.Version = "1.0"
	app.Usage = "thomas: a famous little tank engine run job from fat_controller"

	app.Action = runApp

	app.Flags = []cli.Flag{
		grpcPort,
		standalone,
		fatControllerAddr,
	}
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

	if ctx.String("addr") == "self" {
		return cli.Exit("run failed", 1)
	}

	return nil
}

func exit(sig os.Signal) {
	fmt.Printf("recv %s signal to exit\n", sig)
}
