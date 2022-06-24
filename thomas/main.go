package main

import (
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"path/filepath"
	"sodor/base"
	"sort"
)

// todo initSignals
func main() {
	app := NewApp()

	app.Action = func(ctx *cli.Context) error {
		if ctx.String("addr") == "self" {
			return cli.Exit("run failed", 1)
		}
		return nil
	}

	app.Flags = []cli.Flag{
		grpcPort,
		standalone,
		fatControllerAddr,
	}
	app.Commands = []*cli.Command{
		supervisorCommand,
		httpCommand,
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

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.UseShortOptionHandling = true
	app.Name = filepath.Base(os.Args[0])
	app.Version = "1.0"
	app.Usage = "thomas: a famous little tank engine run task from fat_controller"

	return app
}
