// Pulled and modified from Butter's original implementation of Information Retrieval
// (https://github.com/a-shine/butter/blob/main/retrieve/retrieve.go commit 20ffb299fb196bfe0386ee8ab02987b0fc5e0119)

package pcg

import (
	"errors"
	"fmt"

	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/utils"
)

// retrieve behaviour for a PCG node. When queried, it will either return the information if it is part of the group
// responsible for hosting it, else it will return its known hosts so that the querying node can continue querying the
// network.
func retrieve(overlay node.Overlay, query []byte) []byte {
	persistOverlay := overlay.(*Peer)
	group, err := persistOverlay.Group(string(query))
	if err == nil {
		return append([]byte("found/"), group.Data[:]...)
	}

	//Otherwise not found, return byte array of known hosts to allow for further search...
	hostsStruct := persistOverlay.Node().KnownHostsStruct()
	knownHostsJson := hostsStruct.JsonDigest() // TODO: need to fix this
	return append([]byte("try/"), knownHostsJson...)
}

// AppendRetrieveBehaviour to the Butter node (much like registering an http route in a tradition backend web framework)
func AppendRetrieveBehaviour(node *node.Node) {
	node.RegisterServerBehaviour("pcgRetrieve/", retrieve)
}

// NaiveRetrieve entrypoint to search for a specific piece of information on the network by UUID (information hash)
func NaiveRetrieve(overlay *Peer, query string) ([]byte, error) {
	// Look if I have the information, else query known hosts for information
	// One query per piece of information (one-to-one) hence the query has to be unique i.e i.d.

	// do I have this information, if so return it
	// else BFS (pass the query on to all known hosts (partial view)
	fmt.Println(query)
	block, err := overlay.Group(query)
	if err == nil {
		return block.Data[:], nil
	}
	return bfs(overlay, query)
}

//PROBLEMS::
// There are a series of (potential) problems with the following function:
// 		* we're using len(queue) whilst also updating (this may or may not be a problem depending on how go works)
// 		* Need to make sure no node already queried is checked again, this is by removing duplicates from the queue,
///			 and also having a second list storing which have been checked so they're not added again either

// bfs across the network until information is found. This is not particularly well suited to production and won't scale
// well. However, for testing it provides a deterministic means of checking if information exists on the network.
func bfs(overlay *Peer, query string) ([]byte, error) {
	// Initialise an empty queue
	queue := make([]utils.SocketAddr, 0)
	// Add all my known hosts to the queue
	for host := range overlay.Node().KnownHosts() {
		//print("\nhost", host)
		queue = append(queue, host)
	}
	print("queue: ", len(queue))
	// host map for checked check
	// iterate through knew know hosts
	// only add to queue if not already checked

	for len(queue) > 0 { //TODO CHECK THIS this with go
		// Pop the first element from the queue
		host := queue[0]
		queue = queue[1:]
		// Start a connection to the host, Ask host if he has data, receive resposnse
		response, err := utils.Request(host, []byte("pcgRetrieve/"), []byte(query))
		if err != nil {
			fmt.Println("error in request")
			fmt.Println(err)
			//fmt.Println(response)
		}
		route, payload, err := utils.ParsePacket(response)
		if err != nil {
			fmt.Println("unable to parse packet")
			fmt.Println(err)
			//fmt.Println(response)
		}
		// If the returned packet is success + the data then return it
		// else add the known hosts of the remote node to the end of the queue
		if string(route) == "found/" {
			return payload, nil
		}
		// failed but gave us their known hosts to add to queue
		remoteKnownHosts, _ := utils.AddrSliceFromJson(payload)
		print(payload)
		queue = append(queue, remoteKnownHosts...) // add the remote hosts to the end of the queue. Why does this not loop forever??? GOing in circles innit, may be because len(queue) worked out before and not updated
	}
	return []byte(""), errors.New("failed to retrieve information")
}
