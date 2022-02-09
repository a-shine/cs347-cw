package pcg

import (
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/utils"
	"time"
)

// --- Naive store is included for the sake of comparison ---

// NaiveStore stores information on the network naively by simply placing it on
// the local node. It generate a UUIS for the information and creates an
// information block and return information uuid
func NaiveStore(overlay OverlayPCG, data string) string {
	uuid := overlay.AddGroup(data)
	return uuid
}

// ----------------------------------------------------------

func status(overlayInterface node.Overlay, groupId []byte) []byte {
	overlay := overlayInterface.(*OverlayPCG)
	_, err := overlay.Group(string(groupId))
	if err != nil {
		return []byte("Group not found")
	}
	return []byte("Group member")
}

func storageStatus(overlayInterface node.Overlay, _ []byte) []byte {
	overlay := overlayInterface.(*OverlayPCG)
	// if len(node.Groups()) vs cap(node.Groups()) if len == cap the unable to
	// store more groups if len < cap the able to store more groups
	if len(overlay.Groups()) < cap(overlay.Groups()) {
		return []byte("storage available")
	}
	// if len > cap this should never happen - we should not use more memory
	// than we have allocated to the node at runtime
	return []byte("can't join group")
}

func AppendGroupStoreBehaviour(node *node.Node) {
	node.RegisterServerBehaviour("in-group?/", status)
	node.RegisterServerBehaviour("have-storage?/", storageStatus)

	// much like you can have several server beahbiour you should be able to
	// append several of your own client behaviours to a node (each client
	// behaviour is append to a list and upon startup we create a goroutine for
	// each registered client behaviour)
	node.RegisterClientBehaviour(heartbeat)
}

// NaiveGroupStore is a more intersting implementation to solve the persistent
// storage problem on a high churn network. When we store information we create
// a group that will store this information, and then the participants are
// responsible for filling out the group (to it's max capacity of 3) - groups
// want to be in a state where they have 3 participants
func NaiveGroupStore(pcg OverlayPCG, data string) string {
	uuid := pcg.AddGroup(data)
	newGroup, _ := pcg.Group(uuid)
	newGroup.AddParticipant(pcg.Node().SocketAddr())
	return uuid
}

// Each node is reposnible for managing his own list of group participants - as
// long as it is done faily effectively this should be good enough (no need for
// concensus - should naurally come to concenus as each node manages it's own
// participant list
func heartbeat(overlayInterface node.Overlay) {
	overlay := overlayInterface.(*OverlayPCG)
	for {
		manageParticipants(overlay)
		time.Sleep(time.Second * 30)
	}
}

func manageParticipants(pcg *OverlayPCG) {
	for id, group := range pcg.Groups() { // for all my groups
		// check status of each participant in group
		for _, participant := range group.Participants() {
			// if participant is not alive
			repsonce, err := utils.Request(participant, []byte("group-status/"), []byte(id.String()))
			// remove participant
			if err != nil || string(repsonce) != "Group not found" {
				group.RemoveParticipant(participant)
			}
			// if in group our list of participants is correct
		}
		if len(group.Participants()) < 3 {
			go findParticipants(pcg, &group) // group is in a fragile unhappy state - find more participants
		}
	}
}

func findParticipants(pcg *OverlayPCG, group *Group) {
	for { // runs until a partipant is found - then breaks out of loop
		for _, host := range pcg.Node().KnownHosts() {
			// ask if they would like to join the group i.e. if they have capacity
			repsones, err := utils.Request(host, []byte("storage-status/"), nil)
			if err != nil || string(repsones) != "no storage available" {
				// to bad
			}
			if string(repsones) == "storage available" {
				group.AddParticipant(host)
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

		time.Sleep(time.Second * 30)
	}

}
