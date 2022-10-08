package task_runner

import (
	"log"
	"os"
)

var (
	Warn *log.Logger
	Info *log.Logger
)

func init() {
	Warn = log.New(os.Stderr, "[Task_Runner] ", log.Lshortfile)
	Info = log.New(os.Stdout, "[Task_Runner] ", log.Lshortfile)
}
