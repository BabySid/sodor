package task_runner

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"os"
	"os/exec"
	"syscall"
)

func createCmd(dir string) *cmd.Cmd {
	opt := cmd.Options{
		Buffered:   false,
		Streaming:  false,
		BeforeExec: []func(cmd *exec.Cmd){setCmdDetached},
	}

	script := fmt.Sprintf("%s ----task.identity=run_task &>log", os.Args[0])
	c := cmd.NewCmdOptions(opt, "bash", "-c", script)
	c.Dir = dir
	return c
}

func setCmdDetached(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{}
	///setDeathSig(c.SysProcAttr)
}

//func setDeathSig(sysProcAttr *syscall.SysProcAttr) {
//	sysProcAttr.Setpgid = true
//	sysProcAttr.Pdeathsig = syscall.SIGTERM
//}
