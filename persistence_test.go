package pcg

import (
	"github.com/a-shine/butter/node"
	"testing"
)

//func doNothing(node *node.Node) {
//	for {
//		time.Sleep(5 * time.Second)
//	}
//}
//
//func simple_demo(t *testing.T) {
//	butterNode, _ := node.NewNode(0, 512, doNothing, false)
//	butter.Spawn(&butterNode, false)
//}
//
//func test_persistence(t *testing.T) {
//	// add random data to store
//	// artificially remove and re-add nodes
//	// see if I can still find the data
//}

func simple_test(t *testing.T) {
	butterNode, _ := node.NewNode(0, 512, false)
	SpawnPCG(&butterNode)
}
