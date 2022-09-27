package scheduler

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type dag struct {
	g         *simple.DirectedGraph
	topoNodes []graph.Node
}

func newDAG() *dag {
	return &dag{g: simple.NewDirectedGraph()}
}

func (d *dag) buildFromJob(job *sodor.Job) error {
	for _, t := range job.Tasks {
		d.g.AddNode(simple.Node(t.Id))
	}

	for _, r := range job.Relations {
		d.g.SetEdge(d.g.NewEdge(simple.Node(findTaskIDByName(job.Tasks, r.FromTask)), simple.Node(findTaskIDByName(job.Tasks, r.ToTask))))
	}

	var err error
	d.topoNodes, err = topo.Sort(d.g)
	if err != nil {
		return err
	}

	return nil
}

func findTaskIDByName(ts []*sodor.Task, name string) int32 {
	for _, t := range ts {
		if t.Name == name {
			return t.Id
		}
	}
	gobase.AssertHere()
	return 0
}

func findTaskByID(ts []*sodor.Task, id int32) *sodor.Task {
	for _, t := range ts {
		if t.Id == id {
			return t
		}
	}
	gobase.AssertHere()
	return nil
}
