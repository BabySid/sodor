package jsonrpc

import (
	"errors"
	"fmt"
	"github.com/BabySid/proto/sodor"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
	"sodor/fat_controller/metastore"
	"strings"
)

func checkTaskValid(job *sodor.Job, create bool) error {
	if create && job.Id != 0 {
		return errors.New("job.id must not be set")
	}
	if !create && job.Id == 0 {
		return errors.New("job.id must be set")
	}

	if len(job.Name) >= metastore.MaxNameLen {
		return fmt.Errorf("job.name is long than %d", metastore.MaxNameLen)
	}

	if len(strings.TrimSpace(job.Name)) == 0 {
		return errors.New("job.name is empty")
	}

	if len(job.GetTasks()) == 0 {
		return errors.New("empty tasks")
	}

	s := make(map[string]int)
	taskID := 0
	for _, task := range job.GetTasks() {
		if create && task.Id > 0 {
			return fmt.Errorf("task.id must not be set")
		}
		if len(task.Name) >= metastore.MaxNameLen {
			return fmt.Errorf("task.name is long than %d", metastore.MaxNameLen)
		}

		if len(strings.TrimSpace(task.Name)) == 0 {
			return fmt.Errorf("task.name is empty")
		}

		if len(strings.TrimSpace(task.Script)) == 0 {
			return fmt.Errorf("task.script is empty")
		}

		if task.SchedulerMode != sodor.SchedulerMode_SM_None && task.GetRoutineSpec() != nil && task.GetRoutineSpec().CtSpec == "" {
			return fmt.Errorf("task.spec must be set")
		}

		if _, ok := s[task.Name]; ok {
			return fmt.Errorf("task.name is duplicated in the job")
		}
		s[task.Name] = taskID
		taskID++
	}

	for _, rel := range job.GetRelations() {
		if _, ok := s[rel.FromTask]; !ok {
			return fmt.Errorf("from_task in relations is not exist")
		}

		if _, ok := s[rel.ToTask]; !ok {
			return fmt.Errorf("to_task in relations is not exist")
		}
	}

	// check cycles
	g := simple.NewDirectedGraph()
	for _, tid := range s {
		g.AddNode(simple.Node(tid))
	}
	for _, rel := range job.GetRelations() {
		g.SetEdge(g.NewEdge(simple.Node(s[rel.FromTask]), simple.Node(s[rel.ToTask])))
	}

	cycles := topo.DirectedCyclesIn(g)
	if len(cycles) != 0 {
		return fmt.Errorf("there are cycles in the job")
	}

	return nil
}
