package base

import (
	"github.com/BabySid/gobase"
	"github.com/urfave/cli/v2"
)

// HelpData is a one shot struct to pass to the usage template
type HelpData struct {
	App        interface{}
	FlagGroups []FlagGroup
}

// FlagGroup is a collection of flags belonging to a single topic.
type FlagGroup struct {
	Name  string
	Flags []cli.Flag
}

var (
	LocalHost string
	// AppHelpTemplate is the test template for the default, global app help topic.
	AppHelpTemplate = `NAME:
   {{.App.Name}} - {{.App.Usage}}

   Copyright 2022 The BabySid Authors

USAGE:
   {{.App.HelpName}} [options]{{if .App.Commands}} [command] [command options]{{end}} {{if .App.ArgsUsage}}{{.App.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if .App.Version}}
VERSION:
   {{.App.Version}}
   {{end}}{{if len .App.Authors}}
AUTHOR(S):
   {{range .App.Authors}}{{ . }}{{end}}
   {{end}}{{if .App.Commands}}
COMMANDS:
   {{range .App.Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}{{if .FlagGroups}}
{{range .FlagGroups}}{{.Name}} OPTIONS:
  {{range .Flags}}{{.}}
  {{end}}
{{end}}{{end}}{{if .App.Copyright }}
COPYRIGHT:
   {{.App.Copyright}}
   {{end}}
`
)

func init() {
	ip, err := gobase.GetLocalIP()
	gobase.True(err == nil)
	LocalHost = ip
}
