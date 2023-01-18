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

		if len(strings.TrimSpace(task.Content)) == 0 {
			return fmt.Errorf("task.script is empty")
		}

		if len(task.RunningHosts) == 0 {
			return fmt.Errorf("task.running_hosts must have at least one host")
		}

		if task.JobId != job.Id {
			return fmt.Errorf("task.job_id must be equal with job.id")
		}

		for _, host := range task.RunningHosts {
			if host.Type != sodor.HostType_HostType_IP || host.Node == "" {
				return fmt.Errorf("task.running_hosts.item must be IP")
			}
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

func checkAlertGroupValid(alert *sodor.AlertGroup, create bool) error {
	if create && alert.Id != 0 {
		return errors.New("alert_group.id must not be set")
	}
	if !create && alert.Id == 0 {
		return errors.New("alert_group.id must be set")
	}

	if len(alert.Name) >= metastore.MaxNameLen {
		return fmt.Errorf("alert.name is long than %d", metastore.MaxNameLen)
	}

	if len(strings.TrimSpace(alert.Name)) == 0 {
		return fmt.Errorf("alert.name is empty")
	}

	if len(alert.PluginInstances) == 0 {
		return fmt.Errorf("alert.plugin_instances is empty")
	}

	return nil
}

func checkAlertPluginValid(plugin *sodor.AlertPluginInstance, create bool) error {
	if create && plugin.Id != 0 {
		return errors.New("plugin.id must not be set")
	}
	if !create && plugin.Id == 0 {
		return errors.New("plugin.id must be set")
	}

	if len(plugin.Name) >= metastore.MaxNameLen {
		return fmt.Errorf("plugin.name is long than %d", metastore.MaxNameLen)
	}

	if len(strings.TrimSpace(plugin.Name)) == 0 {
		return fmt.Errorf("plugin.name is empty")
	}

	if plugin.PluginName != sodor.AlertPluginName_APN_DingDing.String() {
		return fmt.Errorf("plugin.plugin_name is invalid")
	}

	if plugin.Dingding == nil {
		return fmt.Errorf("plugin.dingding must be set")
	}

	return nil
}
