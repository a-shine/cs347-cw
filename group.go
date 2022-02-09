package pcg

import (
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/utils"
	"unsafe"
)

const GroupSize = uint64(unsafe.Sizeof(Group{}))

// A Group is the atomic unit of storage. Extra meta-data is attached to the data field (i.e. keywords,
// part-count and geotag) to improve storage and IR performance. A Group is uniquely identified by combining its uuid
//	and part number e.g. <UUID>/<PartNumber>.
type Group struct {
	participants [3]utils.SocketAddr
	part         uint64
	parts        uint64
	data         [4096]byte
}

// Required to make OverlayPCG compatible with the butter Overlay interface
func (o OverlayPCG) Node() *node.Node {
	return o.Node()
}

// Data field getter for a Group
func (g *Group) Data() []byte {
	return g.data[:]
}

func (g *Group) Participants() [3]utils.SocketAddr {
	return g.participants
}

// AddParticipant to Group
func (g *Group) AddParticipant(host utils.SocketAddr) {

}

func (g *Group) RemoveParticipant(host utils.SocketAddr) {

}
