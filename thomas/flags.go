package main

import (
	"github.com/urfave/cli/v2"
	"sodor/base"
)

var (
	grpcPort = &cli.IntFlag{
		Name:        "grpc.port",
		DefaultText: "random",
		Usage:       "GRPC server listening port",
		Value:       0,
	}

	standalone = &cli.BoolFlag{
		Name:  "standalone",
		Usage: "Run as a standalone service, DO NOT interact with fat_controller",
		Value: false,
	}

	fatControllerAddr = &cli.StringFlag{
		Name:  "fat_controller_addr",
		Usage: "Set the fat_controller addresses separated with commas(,)",
	}

	appHelpFlagGroups = []base.FlagGroup{
		{
			Name: "GLOBAL",
			Flags: []cli.Flag{
				grpcPort,
				standalone,
				fatControllerAddr,
			},
		},
	}
)
