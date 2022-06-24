package supervisor

type Supervisor struct {
	// todo event channel
	procMgr *ProcessManager
}

func (s *Supervisor) Init() error {
	return nil
}

func (s *Supervisor) StartProcess() error {
	return nil
}

func (s *Supervisor) StopProcess() error {
	return nil
}

func (s *Supervisor) StopAllProcesses() error {
	return nil
}
