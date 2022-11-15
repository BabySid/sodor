package config

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
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
		Destination: &GetInstance().DataPath,
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

	ConfFile = &cli.StringFlag{
		Name:        "config",
		Usage:       "TOML configuration file",
		DefaultText: "",
		Value:       "",
	}

	DebugMode = &cli.BoolFlag{
		Name:        "debug",
		Usage:       "Set the debug mode",
		DefaultText: "false",
		Value:       false,
	}

	TaskIdentity = &cli.StringFlag{
		Name:        "task.identity",
		Usage:       "Set the identity for running task. e.g. log prefix",
		DefaultText: "RunTask",
		Value:       "RunTask",
		Destination: &GetInstance().TaskIdentity,
	}

	GlobalFlags = []cli.Flag{
		altsrc.NewStringFlag(ListenAddr),
		altsrc.NewStringFlag(DataPath),
		altsrc.NewStringFlag(LogLevel),
		altsrc.NewStringFlag(LogPath),
		altsrc.NewIntFlag(LogMaxAge),
		altsrc.NewBoolFlag(DebugMode),

		ConfFile,
	}

	AppHelpFlagGroups = []base.FlagGroup{
		{
			Name:  "GLOBAL",
			Flags: GlobalFlags,
		},
	}
)
