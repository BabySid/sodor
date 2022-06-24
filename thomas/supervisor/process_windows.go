//go:build windows

package supervisor

import (
	"os"
	"syscall"
)

func setDeathSig(_ *syscall.SysProcAttr) {
}

// Stop stops the command by sending its process group a SIGTERM signal.
// Stop is idempotent. An error should only be returned in the rare case that
// Stop is called immediately after the command ends but before Start can
// update its internal state.
func terminateProcess(pid int) error {
	p := &os.Process{Pid: pid}
	return p.Kill()
}

func setUserID(_ *syscall.SysProcAttr, _ uint32, _ uint32) {
}

func processIsRunning(p *os.Process) bool {
	proc, err := os.FindProcess(p.Pid)
	return proc != nil && err == nil
}
