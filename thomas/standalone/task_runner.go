package standalone

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
)

type TaskRunner interface {
	Run()
}

func GetRunner(typ sodor.TaskType) TaskRunner {
	if typ == sodor.TaskType_TT_Shell {
		return &ShellRunner{}
	}

	gobase.AssertHere()
	return nil
}
