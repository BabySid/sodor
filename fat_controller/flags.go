package main

import (
	"github.com/urfave/cli/v2"
	"sodor/base"
)

var (
	listenAddr = &cli.StringFlag{
		Name:        "listen_addr",
		Usage:       "Set the listen address",
		DefaultText: ":9527",
		Value:       ":9527",
	}

	metaStore = &cli.StringFlag{
		Name:     "metastore.addr",
		Usage:    "Set the metastore address",
		Required: true,
	}

	initMetaStore = &cli.BoolFlag{
		Name:        "init_metastore",
		Usage:       "Init the metastore database. e.g. initialize the tables. this flag is used for one-time operation",
		DefaultText: "false",
		Value:       false,
	}

	logLevel = &cli.StringFlag{
		Name:        "log.level",
		Usage:       "Set the log level",
		DefaultText: "info",
		Value:       "info",
	}

	logPath = &cli.StringFlag{
		Name:        "log.path",
		Usage:       "Set the path for writing the log",
		DefaultText: ".",
		Value:       ".",
	}

	logMaxAge = &cli.IntFlag{
		Name:        "log.max_age",
		Usage:       "Set the max age for the log file ",
		DefaultText: "24*7 hours",
		Value:       24 * 7,
	}

	debugMode = &cli.BoolFlag{
		Name:        "debug",
		Usage:       "Set the debug mode",
		DefaultText: "false",
		Value:       false,
	}

	appHelpFlagGroups = []base.FlagGroup{
		{
			Name: "GLOBAL",
			Flags: []cli.Flag{
				listenAddr,
				metaStore,
				initMetaStore,
				logLevel,
				logPath,
				logMaxAge,
				debugMode,
			},
		},
	}
)
