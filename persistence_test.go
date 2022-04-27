package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/butter-network/pcg-overlay/pcg"
)

// --- Parameters for the test ---
// Changing the following parameters will allow you to test the network in different scenarios

const storerCount = 100 // Number of nodes created with the sole purpose of inserting information into the network

const lifetime = 0      // How long nodes live for (in seconds)
const chanceToDie = 100 // 1 in x chance to die every second
const gracePeriod = 10  // Time node has after spawning before it can die i.e. how long before they can begin to have a chance to die
const churnTime = 60    // Number of seconds the network should churn for
const settleTime = 90   // Amount of time you wish to give the network to settle after churn

// const requestRate = 1
// const requesterCount = 1
// const dataGenInterval = 10 //seconds
// const dataSize = 100

// Global variables to store test data
var activeStorers = 0
var requests = 0
var successRequests = 0
var failedRequests = 0
var storedData = make([]string, 0)
var active = true
var initi = true
var finished = false
var churn = true

// Testing that PCG works on a network with no churn
//func TestNoChurn(t *testing.T) {
//	churn = false
//	fmt.Println("--- Testing the network without churn ---")
//	go maintainNodes()                   // Create the network
//	time.Sleep(settleTime * time.Second) // Let the network settle
//	go makeRequester()                   // Create the requester to check for data persistence
//	fmt.Println("Waiting for requests to finish...")
//	fmt.Println("Active storers (nodes):", activeStorers)
//	// Wait for the requests to have been made
//	for {
//		if finished {
//			break
//		}
//		time.Sleep(1 * time.Second)
//	}
//	fmt.Printf("\n\ntried: %d, failed: %d, len of data %d\n", requests, failedRequests, len(storedData))
//	fmt.Printf("\npercent success: %f\n", (float64(successRequests)/float64(requests))*100)
//}

// 'Churney' network test
// 1. Spawn nodes
// 2. Let them settle
// 3. Churn the network
// 4. Let them settle
// 5. Query the network and see what % of queries are successful
func TestPostChurn(t *testing.T) {
	fmt.Println("--- Testing the network with churn introduced ---")
	go maintainNodes()                   // Create the network
	time.Sleep(churnTime * time.Second)  // Leave the network to churn
	churn = false                        // Stop the network churning
	time.Sleep(settleTime * time.Second) // Let the network settle
	go makeRequester()                   // Create the requester to check for data persistence
	fmt.Println("Waiting for requests to finish...")
	fmt.Println("Active storers (nodes):", activeStorers)
	// Wait for the requests to have been made
	for {
		if finished {
			break
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\n\ntried: %d, failed: %d, len of data %d\n", requests, failedRequests, len(storedData))
	fmt.Printf("\npercent success: %f\n", (float64(successRequests)/float64(requests))*100)
}

// maintainNodes ensures the correct number of nodes are always active
func maintainNodes() {
	for {
		if !active {
			return
		}
		if activeStorers < storerCount {
			var z = storerCount - activeStorers
			for i := 0; i < z; i++ {

				activeStorers = activeStorers + 1
				go makeStorer(initi)
			}
		}
		initi = false
	}
}

// makeStorer creates a node dedicated to storing data. The createData flag dictates whether this node should create
// its own data or not.
func makeStorer(createData bool) {
	butterNode, _ := node.NewNode(0, 512)
	if createData {
		butterNode.RegisterClientBehaviour(addRandomData)
	}
	if lifetime != 0 {
		butterNode.RegisterClientBehaviour(dieAfterX)
	}
	//enable to test churn
	if chanceToDie != 0 {
		butterNode.RegisterClientBehaviour(randomDeath)
	}

	overlay := pcg.NewPCG(butterNode, 512) // Creates a new overlay network
	pcg.AppendRetrieveBehaviour(overlay.Node())
	pcg.AppendGroupStoreBehaviour(overlay.Node())

	butter.Spawn(&overlay, false) // blocking
}

// makeRequester creates a node dedicated to requesting data from the network
func makeRequester() {
	butterNode, _ := node.NewNode(0, 512)
	butterNode.RegisterClientBehaviour(checkPersistence)
	butterNode.RegisterClientBehaviour(dieAfterX)
	overlay := pcg.NewPCG(butterNode, 512) // Creates a new overlay network
	pcg.AppendRetrieveBehaviour(overlay.Node())
	pcg.AppendGroupStoreBehaviour(overlay.Node())

	butter.Spawn(&overlay, false) // blocking
}

// addRandomData generates some data and stores through a node
func addRandomData(overlayInterface node.Overlay) {
	time.Sleep(1 * time.Second)
	peer := overlayInterface.(*pcg.Peer)

	uuid := pcg.Store(peer, gofakeit.Name())
	storedData = append(storedData, uuid)
}

// dieAfterX kills the node after a set amount of time has passed
func dieAfterX(overlayInterface node.Overlay) {
	time.Sleep(time.Duration(lifetime) * time.Second)
	activeStorers = activeStorers - 1
	overlayInterface.Node().Shutdown()
}

// checkPersistence checks if stored data has persisted on the overlay network
func checkPersistence(overlayInterface node.Overlay) {
	peer := overlayInterface.(*pcg.Peer)
	fmt.Println("Retreiver address: ", peer.Node().SocketAddr())

	for {
		if len(peer.Node().KnownHosts()) > 0 { // wait for the requester to connect to a node in the netwrork
			fmt.Println(peer.Node().KnownHosts())
			fmt.Println("Connecting to Network...")
			break
		}
	}

	for _, data := range storedData {
		_, err := pcg.NaiveRetrieve(peer, data)
		//fmt.Println("Retrieved: ", retrieved)
		if err != nil {
			failedRequests = failedRequests + 1
		} else {
			successRequests = successRequests + 1
		}
		requests = requests + 1
	}
	finished = true
}

// randomDeath is a behaviour that can cause a node to fail with a certain chance every second
func randomDeath(overlayInterface node.Overlay) {
	time.Sleep(gracePeriod * time.Second)
	for {
		if !churn {
			return
		}
		num := gofakeit.Number(0, chanceToDie)
		if num == 0 {
			activeStorers = activeStorers - 1
			overlayInterface.Node().Shutdown()
			return
		}
		time.Sleep(1 * time.Second)
	}
}
