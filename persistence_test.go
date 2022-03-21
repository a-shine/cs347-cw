package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/cs347-cw/pcg"
	"github.com/brianvoe/gofakeit/v6"
)

const storerCount = 100

//const requesterCount = 1

const lifetime = 0 // seconds - how long nodes live for
//any of below can be 0
const chanceToDie = 100 // 1 in x chance to die every second
const gracePeriod = 10  //The time nodes have aafter spawning before they can die - how long before they can begin to have a chance to die
const churnTime = 60    //The number of seconds the network should churn for
const settleTime = 90   //The amount of time you wish to give the network to settle after churn, can be 0

// const requestRate = 1 //

// const dataGenInterval = 10 //seconds

// const dataSize = 100

var activeStorers = 0
var requests = 0
var successRequests = 0
var failedRequests = 0
var storedData = make([]string, 0)
var active = true
var initi = true
var finished = false
var churn = true

// Create n nodes and let n/2 nodes exit the network gracefully
func TestNoFailure(t *testing.T) {
	go maintainNodes()                   // create the network
	time.Sleep(churnTime * time.Second)  // leave the network to churn
	churn = false                        //stop the network churning
	time.Sleep(settleTime * time.Second) // let the network settle
	go makeRequester()                   // create the requester to check for data persistence
	fmt.Println("waiting for requests to fininsh, active storers: ", activeStorers)
	for { //wait for the requests to have been made
		if finished {
			break
		}
		time.Sleep(1 * time.Second)
	}
	// time.Sleep(100 * time.Second)
	fmt.Printf("\n\ntried: %d, failed: %d, len of data %d\n", requests, failedRequests, len(storedData))
	fmt.Printf("\npercent success: %f\n", (float64(successRequests)/float64(requests))*100)
}
func TestSomethingElse(t *testing.T) {

}

/*
* Ensures the correct number of nodes are always active
 */

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

/*
* Creates a node dedicated to storing data
* createData dictates whether this node should create it's own data or not
*
 */

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

/*
* Makes a node that will check for all added data
 */
func makeRequester() {
	butterNode, _ := node.NewNode(0, 512)
	butterNode.RegisterClientBehaviour(checkPersistence)
	butterNode.RegisterClientBehaviour(dieAfterX)
	overlay := pcg.NewPCG(butterNode, 512) // Creates a new overlay network
	pcg.AppendRetrieveBehaviour(overlay.Node())
	pcg.AppendGroupStoreBehaviour(overlay.Node())

	butter.Spawn(&overlay, false) // blocking
}

/*
* Generates some data and stores it in the node
 */
func addRandomData(overlayInterface node.Overlay) {
	time.Sleep(1 * time.Second)
	peer := overlayInterface.(*pcg.Peer)

	uuid := pcg.Store(peer, gofakeit.Name())
	storedData = append(storedData, uuid)
}

/*
* Kills the node after a set amount of time has passed
 */
func dieAfterX(overlayInterface node.Overlay) {
	time.Sleep(time.Duration(lifetime) * time.Second)
	activeStorers = activeStorers - 1
	overlayInterface.Node().Shutdown()

}

/*
* Checks stored data has persisted
 */
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

/*
* has a chance to kill a node ever second
 */
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
