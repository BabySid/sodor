package main

import "github.com/urfave/cli/v2"

var (
	supervisorCommand = &cli.Command{
		Name:        "supervisor",
		Aliases:     nil,
		Usage:       "Supervisor process control system for UNIX",
		Category:    "Task Commands",
		Subcommands: nil,
		Flags:       nil,
	}

	httpCommand = &cli.Command{
		Name:        "http",
		Aliases:     nil,
		Usage:       "HttpTask send a request via Http",
		Category:    "Task Commands",
		Subcommands: nil,
		Flags:       nil,
	}
)
