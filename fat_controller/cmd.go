package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"sodor/fat_controller/metastore"
)

var (
	initMetaStore = cli.Command{
		Action:    runInitMetaStore,
		Name:      "init_metastore",
		Usage:     "Init the metastore database. e.g. initialize the tables. this flag is used for one-time operation",
		ArgsUsage: "",
		Flags:     nil,
		Category:  "",
	}
)

func runInitMetaStore(ctx *cli.Context) error {
	err := metastore.GetInstance().AutoMigrate()
	if err != nil {
		log.Warnf("init meta automigrate failed. err=%s", err)
		return err
	}

	return nil
}
