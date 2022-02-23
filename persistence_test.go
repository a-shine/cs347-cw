package main

import (
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/cs347-cw/pcg"
	"math/rand"
	"testing"
	"time"
)

const nodeCount = 5
const listDataHashes []string

// Create n nodes and let n/2 nodes exit the network gracefully
func TestNoFailure(t *testing.T) {
	for i := 0; i < nodeCount; i++ {
		go createNode()
		// add some random data
	}
	butterNode, _ := node.NewNode(0, 512, false)
	butterNode.RegisterClientBehaviour(interact)

	overlay := pcg.NewPCG(&butterNode, 512) // Creates a new overlay network
	//pcg.AppendRetrieveBehaviour(overlay.Node())
	//pcg.AppendGroupStoreBehaviour(overlay.Node())
	
	butter.Spawn(&overlay, false) // blocking
}

//// Find a way of killing goroutines + creating new nodes on the fly and let n/2 exit gracefully + random ungraceful failures
//func TestWithFailure(t *testing.T) {
//}
//
//// All but 1 node dies all nodes fail ungracefully
//func TotalFailure(t *testing.T) {
//}

func addRandomData(overlayInterface node.Overlay) {
	peer := overlayInterface.(*pcg.Peer)
	fmt.Println("Sock addr: ", peer.Node().SocketAddr())
	pcg.PCGStore(peer, String(100))
}

func createNode() {
	butterNode, _ := node.NewNode(0, 512, false)
	butterNode.RegisterClientBehaviour(addRandomData)

	overlay := pcg.NewPCG(&butterNode, 512) // Creates a new overlay network
	pcg.AppendRetrieveBehaviour(overlay.Node())
	pcg.AppendGroupStoreBehaviour(overlay.Node())

	butter.Spawn(&overlay, false) // blocking
}

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
