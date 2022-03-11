package pcg

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	"github.com/a-shine/butter/utils"
)

const GroupStructSize = unsafe.Sizeof(Group{})
const DataReplicationCount = 3
const ParticipantCount = DataReplicationCount // Number of participants that a group optimises for (alias for data replication count)

// A Group is a collection of nodes and the data that they are responsible for maintaining
type Group struct {
	Participants []utils.SocketAddr
	Data         [4096]byte
}

// -- Constructor ---

func NewGroup(data [4096]byte, participant utils.SocketAddr) *Group {
	return &Group{
		Participants: []utils.SocketAddr{participant},
		Data:         data,
	}
}

// --- Setters ---

func (g *Group) SetParticipants(participants []utils.SocketAddr) {
	g.Participants = participants
}

// --- Add and remove participants ---

// AddParticipant to Group
func (g *Group) AddParticipant(host utils.SocketAddr) error {
	if len(g.Participants) >= ParticipantCount {
		return errors.New("group is full")
	}
	g.SetParticipants(append(g.Participants, host))
	return nil
}

// RemoveParticipant from Group
func (g *Group) RemoveParticipant(host utils.SocketAddr) error {
	fmt.Println("Removing:", host, "from a group")
	for i, participant := range g.Participants {
		fmt.Println(participant.ToString())
		fmt.Println(host.ToString())
		if participant.ToString() == host.ToString() {
			g.Participants = append(g.Participants[:i], g.Participants[i+1:]...)
			break
		}
		if i == len(g.Participants) {
			return errors.New("host not in group")
		}
	}
	return nil
}

// ToJson returns a JSON representation of the group
func (g *Group) ToJson() []byte {
	groupJson, _ := json.Marshal(g)
	return groupJson
}

// String returns a string representation of the group
func (g *Group) String() string {
	return fmt.Sprintf("Data: %s\nGroup Members: %v\nUUID: %x\n\n", g.Data[:], g.Participants, sha256.Sum256(g.Data[:]))
}
