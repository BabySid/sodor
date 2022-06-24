//go:build linux

package supervisor

import (
	"os"
	"syscall"
)

func setDeathSig(sysProcAttr *syscall.SysProcAttr) {
	sysProcAttr.Setpgid = true
	sysProcAttr.Pdeathsig = syscall.SIGKILL
}

func terminateProcess(pid int) error {
	// Signal the process group (-pid), not just the process, so that the process
	// and all its children are signaled. Else, child procs can keep running and
	// keep the stdout/stderr fd open and cause cmd.Wait to hang.
	return syscall.Kill(-pid, syscall.SIGKILL)
}

func setUserID(procAttr *syscall.SysProcAttr, uid uint32, gid uint32) {
	procAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid, NoSetGroups: true}
}

func processIsRunning(p *os.Process) bool {
	return p.Signal(syscall.Signal(0)) == nil
}
