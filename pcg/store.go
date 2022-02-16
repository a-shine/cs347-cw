package pcg

import (
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/utils"
	"time"
	"fmt"
)

const (
	// Storage overlay route API
	inGroupUri = "in-group?/"
	canJoinUri = "can-join?/"
)

func PCGStore(overlay *Peer, data string) string {
	uuid := overlay.AddGroup(data)
	return uuid
}

// AppendGroupStoreBehaviour registers the behaviours that allow the node to work with the pcg overlay
func AppendGroupStoreBehaviour(node *node.Node) {
	node.RegisterServerBehaviour(inGroupUri, inGroup)
	node.RegisterServerBehaviour(canJoinUri, canJoin)

	node.RegisterClientBehaviour(heartbeat)
}

// --- Server behaviours (can be thought of as questions) ---

func inGroup(overlayInterface node.Overlay, groupId []byte) []byte {
	pcg := overlayInterface.(*Peer)
	_, err := pcg.Group(string(groupId))
	if err != nil {
		return []byte("Group not found")
	}
	return []byte("Group member")
}

func canJoin(overlayInterface node.Overlay, payload []byte) []byte {
	peer := overlayInterface.(*Peer)
	// if len(node.Groups()) vs cap(node.Groups()) if len == cap the unable to
	// store more groups if len < cap the able to store more groups
	fmt.Println("I'm can join")
	if peer.currentStorage < peer.maxStorage {
		//Start go routine that will add me to the group that has been requested

		//Parse payload to get the group which I'm supposed to join :)
		

		go 
		return []byte("accept")
	}
	// if len > cap this should never happen - we should not use more memory
	// than we have allocated to the node at runtime
	return []byte("can't join group")
}

// --- Client behaviours ---

// Each node is responsible for managing his own list of group participants - as long as it is done faily effectively
// this should be good enough (no need for concensus - should naurally come to concenus as each node manages it's own
// participant list
func heartbeat(overlayInterface node.Overlay) {
	pcg := overlayInterface.(*Peer)
	for {
		manageParticipants(pcg)
		// fmt.Println("My heart beats for you")
		time.Sleep(time.Second * 5)
	}
}

func manageParticipants(pcg *Peer) {
	for id, group := range pcg.Groups() { // for all my groups
		// check status of each participant in group
		for _, participant := range group.Participants() {
			// if participant is not alive
			repsonce, err := utils.Request(participant, []byte(inGroupUri), id[:])
			// remove participant
			if err != nil || string(repsonce) != "Group not found" {
				group.RemoveParticipant(participant)
			}
			// if in group our list of participants is correct
		}
		print(len(group.Participants()))
		if len(group.Participants()) < 3 { //FIx this as if findParticipants already running then it'll make multiple
			go findParticipants(pcg, &group) // group is in a fragile unhappy state - find more participants
		}
	}
}

func findParticipants(pcg *Peer, group *Group) {
	// fmt.Print("finding!")
	for { // runs until a partipant is found - then breaks out of loop
		for _, host := range pcg.Node().KnownHosts() {
			// ask if they would like to join the group i.e. if they have capacity
			response, err := utils.Request(host, []byte(canJoinUri), nil)
			// fmt.Println(string(response))
			if err != nil || string(response) == "no storage available" {
				// too bad
				fmt.Print("life is sad")
			} 
			
			if string(response) == "accepted" {
				group.AddParticipant(host)

				//Send message to host that we want him to be added to our group
				
				if len(group.Participants()) == 3 {
					break
					return // this instead maybe? no need for second check later
				}
			}
		}
		// if that doesn't work - ask other participants in group if they know
		// someone in their known hosts as a last resort if none of your known hosts
		// + none of the other group partipants have a known host available to join
		// a group - ask if your known hosts' known hosts known anyone - i.e. do
		// breadth first search as a last resort to avoid increasing message
		// complexity

		// do this until a participant is found - so if doen't work first time try
		// again - if group particpants becomes 3 then break
		if len(group.Participants()) == 3 {
			break
		}
		time.Sleep(time.Second * 5)
	}

}
