package pcg

import (
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
)

func SpawnPCG(node *node.Node, memory uint64) {
	// create a node with 0.05 memory
	// define size of pcg storage - have a maximum size

	pcgOverlay := NewPCGOverlay(node)

	// append the necessary node behaviours here
	AppendRetrieveBehaviour(node)
	//AppendStoreBehaviour(node)

	butter.Spawn(&pcgOverlay, false) // not using traverse
}
