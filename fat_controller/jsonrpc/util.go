package jsonrpc

import (
	"errors"
	"fmt"
	"github.com/BabySid/proto/sodor"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
	"sodor/fat_controller/metastore"
	"sodor/fat_controller/scheduler"
	"strings"
)

func checkJobValid(job *sodor.Job, create bool) error {
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

	if job.ScheduleMode == sodor.ScheduleMode_ScheduleMode_Crontab {
		if job.GetRoutineSpec() == nil || job.GetRoutineSpec().CtSpec == "" {
			return fmt.Errorf("task.spec must be set")
		}

		parser := scheduler.NewParser()
		if _, err := parser.Parse(job.GetRoutineSpec().CtSpec); err != nil {
			return fmt.Errorf("invalid task.spec. %s", err)
		}
	}

	if len(job.GetTasks()) == 0 {
		return errors.New("empty tasks")
	}

	s := make(map[string]int)
	for i, task := range job.GetTasks() {
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

		if len(task.RunningHosts) != 1 {
			return fmt.Errorf("task.running_hosts must have only one host")
		}

		if task.RunningHosts[0].Type != sodor.HostType_HostType_IP || task.RunningHosts[0].Node == "" {
			return fmt.Errorf("task.running_hosts.item must be IP")
		}

		if _, ok := s[task.Name]; ok {
			return fmt.Errorf("task.name is duplicated in the job")
		}
		s[task.Name] = i
	}

	if len(job.GetTasks()) > 1 && len(job.GetRelations()) == 0 {
		return fmt.Errorf("task.relations must set when there are more than one tasks")
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

	_, err := topo.Sort(g)
	//cycles := topo.DirectedCyclesIn(g)
	if err != nil {
		return fmt.Errorf("there are cycles in the job")
	}

	return nil
}
