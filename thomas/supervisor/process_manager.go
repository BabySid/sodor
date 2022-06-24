package supervisor

import "sync"

type ProcessManager struct {
	processes sync.Map
}

func NewManager() *ProcessManager {
	return &ProcessManager{processes: sync.Map{}}
}

func (pm *ProcessManager) CreateProcess(name string, cmd string) *Process {
	return nil
}

func (pm *ProcessManager) StartPrograms() {

}

func (pm *ProcessManager) StopPrograms() {

}
