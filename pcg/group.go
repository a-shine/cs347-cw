package pcg

import (
	"errors"
	"github.com/a-shine/butter/utils"
	"unsafe"
)

const GroupSize = unsafe.Sizeof(Group{})
const ReplicationCount = 3
const ParticipantCount = ReplicationCount

type Group struct {
	participants [3]utils.SocketAddr
	data         [4096]byte
}

// Data held by group
func (g *Group) Data() []byte {
	return g.data[:]
}

// Participants in group
func (g *Group) Participants() [3]utils.SocketAddr {
	return g.participants
}

// AddParticipant to Group
func (g *Group) AddParticipant(host utils.SocketAddr) error {
	for _, participant := range g.participants {
		if participant.IsEmpty() {
			participant = host
			return nil
		}
	}
	return errors.New("group is full")
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
		participants: [3]utils.SocketAddr{participant, utils.SocketAddr{}, utils.SocketAddr{}},
		data:         data,
	}
}
