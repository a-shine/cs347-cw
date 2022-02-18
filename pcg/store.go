package pcg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/utils"
)

const (
	// Storage overlay route API
	inGroupUri = "in-group?/"
	canJoinUri = "can-join?/"
)

func PCGStore(overlay *Peer, data string) string {
	uuid := overlay.NewGroup(data)
	return uuid
}

var alreadyFinding bool

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
	// fmt.Println(string(payload))
	// if len(node.Groups()) vs cap(node.Groups()) if len == cap the unable to
	// store more groups if len < cap the able to store more groups
	// fmt.Println("I'm can join")
	if peer.currentStorage < peer.maxStorage {
		//Start go routine that will add me to the group that has been requested

		//Parse payload to get the group which I'm supposed to join :)
		var groupDigest Group //TODO update to group Digest struct once group has been filled out further
		err := json.Unmarshal(payload, &groupDigest)
		if err != nil {
			fmt.Println("error marchallng group")
		}
		// fmt.Println(groupDigest.String())
		peer.JoinGroup(groupDigest)
		// fmt.Println("Joined someones group")
		return []byte("accepted")
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
	pcgn := overlayInterface.(*Peer)
	for {
		manageParticipants(pcgn)
		time.Sleep(time.Second * 5)
	}
}

func (p *Peer) amILeader(g *Group) bool {
	socketAddr := p.Node().SocketAddr()
	socketAddrStr := socketAddr.ToString()

	if !GroupContains(g.Participants, socketAddr) {
		return false
	}

	for _, h := range g.Participants {
		if h.ToString() > socketAddrStr {
			return false
		}
	}
	return true
}

func manageParticipants(peer *Peer) {
	//fmt.Println("Num groups", len(peer.Groups()))

	for id, group := range peer.Groups() { // for all my groups
		// check status of each participant in group
		//fmt.Println(group.Participants)
		for _, participant := range group.Participants {
			// if participant is not alive
			responce, err := utils.Request(participant, []byte(inGroupUri), id[:])
			// remove participant
			if err != nil || string(responce) != "Group not found" {
				group.RemoveParticipant(participant)
				fmt.Println(group.Participants)
			}
			// if in group our list of participants is correct
		}
		if peer.amILeader(group) && ((len(group.Participants)) < 3) && !alreadyFinding { //FIx this as if findParticipants already running then it'll make multiple
			go findParticipants(peer, group) // group is in a fragile unhappy state - find more participants
		}
	}
}

func GroupContains(g []utils.SocketAddr, h utils.SocketAddr) bool {
	for _, a := range g {
		if a.ToString() == h.ToString() {
			return true
		}
	}
	return false
}

func findParticipants(pcg *Peer, group *Group) {
	alreadyFinding = true
	// fmt.Print("finding!")
	for { // runs until a partipant is found - then breaks out of loop
		for _, host := range pcg.Node().KnownHosts() {
			if GroupContains(group.Participants, host) {
				continue
			}
			// ask if they would like to join the group i.e. if they have capacity
			output, err := json.Marshal(group)
			// fmt.Println(string(output))
			if err != nil {
				break
			}
			response, err := utils.Request(host, []byte(canJoinUri), output)
			fmt.Println(string(response))
			if err != nil || string(response) == "no storage available" {
				// too bad
				fmt.Println(err)
				fmt.Println("life is sad")
			}

			if string(response) == "accepted" {
				fmt.Println("adding:", host)
				// fmt.Printf("%p", &group.Participants)
				err := group.AddParticipant(host)
				if err != nil {
					fmt.Println(err)
				}

				//Send message to host that we want him to be added to our group
				if len(group.Participants) == 3 {
					break
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

		time.Sleep(time.Second * 5)
		if len(group.Participants) == 3 {
			break
		}
	}
	alreadyFinding = false
}
