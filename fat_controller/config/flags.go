package config

import (
	"github.com/urfave/cli/v2"
	"sodor/base"
)

var (
	ListenAddr = &cli.StringFlag{
		Name:        "listen_addr",
		Usage:       "Set the listen address",
		DefaultText: ":9527",
		Value:       ":9527",
	}

	MetaStore = &cli.StringFlag{
		Name:        "metastore.addr",
		Usage:       "Set the metastore address",
		Required:    true,
		DefaultText: "mysql://$user:$passwd@tcp($host:$port)/$db?charset=utf8mb4&parseTime=True&loc=Local",
		Value:       "",
	}

	InitMetaStore = &cli.BoolFlag{
		Name:        "init_metastore",
		Usage:       "Init the metastore database. e.g. initialize the tables. this flag is used for one-time operation",
		DefaultText: "false",
		Value:       false,
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

	GlobalFlags = []cli.Flag{
		ListenAddr,
		MetaStore,
		InitMetaStore,
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
