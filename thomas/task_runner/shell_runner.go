package task_runner

import (
	"github.com/go-cmd/cmd"
	"google.golang.org/protobuf/types/known/structpb"
	"regexp"
	"strings"
)

type ShellRunner struct {
	*TaskRunner
	req *regexp.Regexp

	outputVars map[string]interface{}
}

func (s *ShellRunner) Run() error {
	s.TaskRunner = NewTaskRunner()
	s.outputVars = make(map[string]interface{})

	if err := s.SetUp(); err != nil {
		Warn.Printf("task SetUp failed. err = %s", err.Error())
		return err
	}

	Info.Printf("task(%d-%d-%s) begin to run", s.request.JobId, s.request.TaskId, s.request.Task.Name)

	s.response.ExitCode = SystemError

	defer func() {
		Info.Printf("task(%d-%d-%s) run finished", s.request.JobId, s.request.TaskId, s.request.Task.Name)
		if err := s.TearDown(); err != nil {
			Warn.Printf("task(%d-%d-%s) TearDown failed. err = %s", s.request.JobId, s.request.TaskId, s.request.Task.Name, err.Error())
		}
	}()

	// [\\$#]\\{set_value\\(([^)]*)\\)}
	req, err := regexp.Compile("set_value\\(([^)]*)\\)")
	if err != nil {
		Warn.Printf("task(%d-%d-%s) regexp.Compile failed. err = %s", s.request.JobId, s.request.TaskId, s.request.Task.Name, err.Error())
		s.response.ExitMsg = err.Error()
		return err
	}
	s.req = req

	c := cmd.NewCmdOptions(cmd.Options{
		Buffered:   false,
		Streaming:  true,
		BeforeExec: nil,
	}, "bash", "-c", s.request.Task.Script)

	go s.processStdoutStderr(c)

	status := <-c.Start()

	vars, err := structpb.NewStruct(s.outputVars)
	if err != nil {
		Warn.Printf("task(%d-%d-%s) structpb.NewStruct(%+v) failed. err = %s",
			s.request.JobId, s.request.TaskId, s.request.Task.Name, s.outputVars, err.Error())
		return err
	}
	s.response.OutputVars = vars

	s.response.ExitCode = int32(status.Exit)
	s.response.ExitMsg = OKMsg
	if status.Error != nil {
		s.response.ExitMsg = status.Error.Error()
	}

	return nil
}

func (s *ShellRunner) findOutputValue(line string) map[string]interface{} {
	ls := s.req.FindStringSubmatch(line)
	if len(ls) > 1 {
		idx := strings.Index(ls[1], "=")
		if idx > 0 {
			key := strings.TrimSpace(ls[1][:idx])
			value := strings.TrimSpace(ls[1][idx+1:])
			return map[string]interface{}{key: value}
		}
	}

	return nil
}

func (s *ShellRunner) processStdoutStderr(c *cmd.Cmd) {
	for c.Stdout != nil || c.Stderr != nil {
		select {
		case line, open := <-c.Stdout:
			if !open {
				c.Stdout = nil
				continue
			}
			vars := s.findOutputValue(line)
			for k, v := range vars {
				s.outputVars[k] = v
			}
			Info.Println(line)
		case line, open := <-c.Stderr:
			if !open {
				c.Stderr = nil
				continue
			}
			Warn.Println(line)
		}
	}
}
