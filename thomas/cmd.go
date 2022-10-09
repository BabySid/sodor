package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"sodor/thomas/config"
	"sodor/thomas/task_runner"
)

var (
	runTask = cli.Command{
		Action:    runTaskCommand,
		Name:      "run_task",
		Usage:     "Run as a task_runner process. DO NOT interact with fat_controller",
		ArgsUsage: "",
		Flags: []cli.Flag{
			config.TaskIdentity,
		},
		Category: "",
	}
)

func runTaskCommand(ctx *cli.Context) error {
	t := task_runner.GetRunner()
	if err := t.Run(); err != nil {
		log.Fatalf("task_runner run failed. err = %s", err.Error())
	}
	return nil
}
