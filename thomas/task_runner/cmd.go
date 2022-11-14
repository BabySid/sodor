package task_runner

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/go-cmd/cmd"
	"os"
	"os/exec"
)

func createCmd(dir string) *cmd.Cmd {
	opt := cmd.Options{
		Buffered:   false,
		Streaming:  false,
		BeforeExec: []func(cmd *exec.Cmd){gobase.SetChildrenProcessDetached},
	}

	script := fmt.Sprintf("%s ----task.identity=run_task &>log", os.Args[0])
	c := cmd.NewCmdOptions(opt, "bash", "-c", script)
	c.Dir = dir

	return c
}
