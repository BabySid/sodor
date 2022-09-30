package standalone

import (
	"log"
	"os"
)

var (
	Warn *log.Logger
	Info *log.Logger
)

func init() {
	Warn = log.New(os.Stderr, "[Job_Runner] ", log.Lshortfile)
	Info = log.New(os.Stdout, "[Job_Runner] ", log.Lshortfile)
}
