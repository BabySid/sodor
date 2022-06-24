package supervisor

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type State int

const (
	// Stopped the stopped state
	Stopped State = iota

	// Starting the starting state
	Starting = 10

	// Running the running state
	Running = 20

	// Backoff the backoff state
	Backoff = 30

	// Stopping the stopping state
	Stopping = 40

	// Exited the Exited state
	Exited = 100

	// Fatal the Fatal state
	Fatal = 200

	// Unknown the unknown state
	Unknown = 1000
)

func (p State) String() string {
	switch p {
	case Stopped:
		return "Stopped"
	case Starting:
		return "Starting"
	case Running:
		return "Running"
	case Backoff:
		return "Backoff"
	case Stopping:
		return "Stopping"
	case Exited:
		return "Exited"
	case Fatal:
		return "Fatal"
	default:
		return "Unknown"
	}
}

type Process struct {
	name      string
	cmdStr    string
	cmd       *exec.Cmd
	startTime time.Time
	stopTime  time.Time
	pid       int
	status    State
	//true if process is starting
	inStart bool

	stopByUser bool

	lock sync.RWMutex
}

func NewProcess(name string, cmd string) *Process {
	p := &Process{
		name:       name,
		cmdStr:     cmd,
		cmd:        nil,
		startTime:  time.Unix(0, 0),
		stopTime:   time.Unix(0, 0),
		pid:        0,
		status:     Stopped,
		inStart:    false,
		stopByUser: false,
	}

	return p
}

func (p *Process) Start() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.inStart {
		return nil
	}

	p.inStart = true
	p.stopByUser = false

	go func() {
		for {
			p.run()

			if p.stopByUser {
				break
			}
		}
		p.inStart = false
	}()

	return nil
}

func (p *Process) run() {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.isRunning() {
		return
	}

	p.startTime = time.Now()

	for !p.stopByUser {
		p.status = Starting

		p.cmd = exec.Command(p.cmdStr)
		p.cmd.SysProcAttr = &syscall.SysProcAttr{}
		setDeathSig(p.cmd.SysProcAttr)
		p.setEnv()

		err := p.cmd.Start()
		if err != nil {
			// todo log it and restart
		}

		p.status = Running

		p.wait()

		if p.status == Running {
			p.status = Exited
			break
		} else {
			p.status = Backoff
		}
	}
}

func (p *Process) wait() {
	p.cmd.Wait()

	if p.cmd.ProcessState != nil {

	}

	p.stopTime = time.Now()
}

func (p *Process) Stop() {
	p.lock.Lock()
	defer p.lock.Unlock()
	
	p.stopByUser = true

	if !p.isRunning() {
		return
	}

	p.sendSignal(syscall.SIGTERM)
}

func (p *Process) GetPid() int {
	// todo check status
	return p.cmd.Process.Pid
}

func (p *Process) isRunning() bool {
	if p.cmd != nil && p.cmd.Process != nil {
		return processIsRunning(p.cmd.Process)
	}
	return false
}

func (p *Process) sendSignal(sig os.Signal) error {
	if p.cmd != nil && p.cmd.Process != nil {
		return terminateProcess(p.cmd.Process.Pid)
	}
	return fmt.Errorf("process is not started")
}

func (p *Process) setEnv() {
	p.cmd.Env = os.Environ()
}
