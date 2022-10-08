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

	TaskRunner = &cli.BoolFlag{
		Name:  "task_runner",
		Usage: "Run as a task_runner process, DO NOT interact with fat_controller",
		Value: false,
	}

	GlobalFlags = []cli.Flag{
		ListenAddr,
		TaskRunner,
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
