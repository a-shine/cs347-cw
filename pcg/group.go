package pcg

import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	"github.com/a-shine/butter/utils"
	// "fmt"
)

const GroupSize = unsafe.Sizeof(Group{})
const ReplicationCount = 3
const ParticipantCount = ReplicationCount

type Group struct {
	Participants []utils.SocketAddr
	Data         [4096]byte
}

func (g *Group) SetParticipants(participants []utils.SocketAddr) {
	g.Participants = participants
}

func (g *Group) ToJson() []byte {
	json, _ := json.Marshal(g)
	return json
}

// AddParticipant to Group
func (g *Group) AddParticipant(host utils.SocketAddr) error {
	// fmt.Printf("%p", &g.Participants)
	if len(g.Participants) >= 3 {
		return errors.New("group is full")
	}
	g.SetParticipants(append(g.Participants, host))
	// fmt.Println(host)
	return nil

}

func removeSockAddrSlice(s []utils.SocketAddr, i int) []utils.SocketAddr {
	s[i] = s[len(s)-1]
	s = s[:len(s)-1]
	fmt.Println("in remov", s)
	return s
}

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

	// fmt.Println("HEREERRERE")
	// g.Participants = removeSockAddrSlice(g.Participants, i)
	// fmt.Println(g.Participants)
	return nil
}

func NewGroup(data [4096]byte, participant utils.SocketAddr) *Group {
	return &Group{
		Participants: []utils.SocketAddr{participant},
		Data:         data,
	}
}

func (g *Group) String() string {
	return string(g.Data[:])
}

func (g *Group) numPart() int {
	return len(g.Participants)
}
