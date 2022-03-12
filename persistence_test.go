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
	pcg.Store(peer, String(100))
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

/*
Testing
 * Other thoughts:
	* What are we doing with nodes that leave the network and join again? Currently we take them as fresh (not holding any data) when they rejoin. Is this what we decide?
	* If nodes are "reinitialised" on rejoining, then we can simulate this easily just with new nodes that contain no data. (only slight difference that I can think of is the underlying butter network known host list churn but that's fine?)

 * General Notes:
	* Don't have to implement all proper functionalilty to every node, such as the ones for just making requests
		they can be seen as nodes outside the network that don't do anythign properly just make requests.
	* How do we want to go about adding data to the network? does each node start with some random amount of data? then don't add any later on?
	* Could make nodes that don't add any, and they're just good helpful bois.

 * Global settings:
	* average lifetime of nodes + additional randomized lifespan + auto kill.
	* % of deaths be graceful or ungraceful.
	* have a delayed communication response to check if nodes are being removed from PCGs if they are only slow not dead?
	* Rate of requests?
	* Number of nodes on the netowrk?
	* Number of nodes making requests? ()
	* vary rate of data addition or max data added?
	* vary mean and s.d. of node capacities
 
 * Functions:
	* Kill myself function:
		* Uses config settings to generate pseudo random lifetime and death style
		* will kill the node when time is up!
	
	* Add data to network
		* Nodes generate random strings of data and send it to other nodes
		* Add UUID (and the data expected for validation purposes) to global array for the sake of retrieve

	* Retreive data
		* The main function for evaluating availability
		* requests data it knows should be on the network
		* report success rate
		* performed at vraying request rates?

	* Network evaluation
		* Everytimme a tcp request is made add to some global counter
		* Used to model network traffic
		* can be used to compare heartbeat to gossip should we choose to do that.

	* Generic Evaluation:
		* work out % availability under different environments/churn rate/scenarios
*/
