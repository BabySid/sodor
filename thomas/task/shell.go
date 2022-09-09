package task

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var (
	scriptFile = &cli.StringFlag{
		Name:     "script_file",
		Usage:    "Set the script file for shell task",
		Required: true,
	}

	ShellCommand = &cli.Command{
		Name:        "shell",
		Aliases:     nil,
		Usage:       "ShellTask runs a shell script",
		Category:    "Task Commands",
		Subcommands: nil,
		Flags: []cli.Flag{
			scriptFile,
		},
		Action: runShell,
	}
)

func runShell(ctx *cli.Context) error {
	fmt.Println("run shell", ctx.String(scriptFile.Name))
	return nil
}
