package config

import (
	"github.com/urfave/cli/v2"
	"sodor/base"
)

var (
	ListenAddr = &cli.StringFlag{
		Name:        "listen_addr",
		Usage:       "Set the listen address",
		DefaultText: ":9528",
		Value:       ":9528",
	}

	DataPath = &cli.StringFlag{
		Name:        "data.path",
		Usage:       "Set the path for the context of tasks. e.g. the task log, response and so on",
		DefaultText: ".",
		Value:       "./",
	}

	LogLevel = &cli.StringFlag{
		Name:        "log.level",
		Usage:       "Set the log level",
		DefaultText: "info",
		Value:       "info",
	}

	LogPath = &cli.StringFlag{
		Name:        "log.path",
		Usage:       "Set the path for writing the log",
		DefaultText: ".",
		Value:       "./",
	}

	LogMaxAge = &cli.IntFlag{
		Name:        "log.max_age",
		Usage:       "Set the max age for the log file ",
		DefaultText: "24*7 hours",
		Value:       24 * 7,
	}

	DebugMode = &cli.BoolFlag{
		Name:        "debug",
		Usage:       "Set the debug mode",
		DefaultText: "false",
		Value:       false,
	}

	RunTask = &cli.BoolFlag{
		Name:  "run_task",
		Usage: "Run as a task_runner process. DO NOT interact with fat_controller",
		Value: false,
	}

	TaskIdentity = &cli.StringFlag{
		Name:        "task.identity",
		Usage:       "Set the identity for running task. e.g. log prefix",
		DefaultText: "RunTask",
		Value:       "RunTask",
	}

	GlobalFlags = []cli.Flag{
		ListenAddr,
		RunTask,
		TaskIdentity,
		DataPath,
		LogLevel,
		LogPath,
		LogMaxAge,
		DebugMode,
	}

	AppHelpFlagGroups = []base.FlagGroup{
		{
			Name:  "GLOBAL",
			Flags: GlobalFlags,
		},
	}
)
