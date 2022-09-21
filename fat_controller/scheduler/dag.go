package scheduler

import (
	"fmt"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type dag struct {
	g *simple.DirectedGraph
}

func newDAG() *dag {
	return &dag{g: simple.NewDirectedGraph()}
}

func (d *dag) buildFromMetaStore() {
	g := simple.NewDirectedGraph()
	g.AddNode(simple.Node(1))
	g.AddNode(simple.Node(2))
	g.SetEdge(g.NewEdge(simple.Node(1), simple.Node(2)))

	cycles := topo.DirectedCyclesIn(g)

	sorted, _ := topo.Sort(g)

	fmt.Println(cycles, sorted)
}

func (d *dag) Cycles() {
	
}
