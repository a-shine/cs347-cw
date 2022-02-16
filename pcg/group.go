package pcg

import (
	"errors"
	"github.com/a-shine/butter/utils"
	"unsafe"
	// "fmt"
)

const GroupSize = unsafe.Sizeof(Group{})
const ReplicationCount = 3
const ParticipantCount = ReplicationCount

type Group struct {
	participants []utils.SocketAddr
	data         [4096]byte
}

// Data held by group
func (g *Group) Data() []byte {
	return g.data[:]
}

// Participants in group
func (g *Group) Participants() []utils.SocketAddr {
	return g.participants
}

func (g *Group) SetParticipants(participants []utils.SocketAddr){
	g.participants=participants

}

// AddParticipant to Group
func (g *Group) AddParticipant(host utils.SocketAddr) error {
	if(len(g.Participants())>=3){
		return errors.New("group is full")
	}
	g.SetParticipants(append(g.Participants(),host))

	return nil
	
}

func (g *Group) RemoveParticipant(host utils.SocketAddr) error {
	for _, participant := range g.participants {
		if participant.ToString() == host.ToString() {
			participant = utils.SocketAddr{}
			return nil
		}
	}
	return errors.New("host not in group")
}

func NewGroup(data [4096]byte, participant utils.SocketAddr) Group {
	return Group{
		participants: []utils.SocketAddr{participant},
		data:         data,
	}
}

func (g *Group) String() string {
	return string(g.data[:])
}
