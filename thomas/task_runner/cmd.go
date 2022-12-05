package task_runner

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/go-cmd/cmd"
	"os"
	"os/exec"
	"path/filepath"
)

func createCmd(dir string) *cmd.Cmd {
	opt := cmd.Options{
		Buffered:   false,
		Streaming:  false,
		BeforeExec: []func(cmd *exec.Cmd){gobase.SetChildrenProcessDetached},
	}

	bin, _ := filepath.Abs(os.Args[0])
	script := fmt.Sprintf("%s run_task --task.identity=%s &>log", bin, dir)
	c := cmd.NewCmdOptions(opt, "bash", "-c", script)
	c.Dir = dir

	return c
}
