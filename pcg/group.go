package pcg

import "github.com/a-shine/butter/utils"

const ReplicationCount = 3
const ParticipantCount = ReplicationCount

type Group struct {
	participants [ParticipantCount]utils.SocketAddr
	data         [4096]byte
}

// Data field getter for a Group
func (b *Group) Data() []byte {
	return b.data[:]
}
