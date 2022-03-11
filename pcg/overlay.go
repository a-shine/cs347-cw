package pcg

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/a-shine/butter/node"
)

// Peer implements an overlay node which is described in the Butter node overlay interface
type Peer struct {
	node           *node.Node
	maxStorage     uint64 //TODO will be in node
	currentStorage uint64 //TODO will be in node
	storage        map[[32]byte]*Group
}

// --- Constructor ---

// NewPCG overlay node
func NewPCG(node *node.Node, maxMemoryMb uint64) Peer {
	maxMemory := MbToBytes(maxMemoryMb)
	maxStorage := MaxStorage(maxMemory)
	fmt.Println("Max storage:", maxStorage)
	return Peer{
		node:       node,
		maxStorage: maxStorage,
		storage:    make(map[[32]byte]*Group),
	}
}

// --- Getters ---

func (o *Peer) Node() *node.Node {
	return o.node
}

func (p *Peer) AvailableStorage() uint64 {
	return p.maxStorage - p.currentStorage
}

// Group from the node's storage by its UUID. If the block is not found, an empty block with an error is returned.
func (o *Peer) Group(id string) (*Group, error) {
	var hash [32]byte
	data, _ := hex.DecodeString(id)
	copy(hash[:], data)
	if group, ok := o.storage[hash]; ok {
		return group, nil
	}
	return nil, errors.New("block not found")
}

func (o *Peer) Groups() map[[32]byte]*Group {
	return o.storage
}

// CreateGroup to the node's storage. A UUID is generated for every bit of information added to the network (no update
// functionality yet!). Returns the UUID of the new block as a string.
func (p *Peer) CreateGroup(data string) string {
	// check if node memory (allocated by the user on initialization) is full - node API
	var formattedData [4096]byte
	copy(formattedData[:], data)
	hsha2 := sha256.Sum256(formattedData[:])
	p.storage[hsha2] = NewGroup(formattedData, p.node.SocketAddr())
	p.currentStorage += 4096 //TODO
	return fmt.Sprintf("%x", hsha2)
}

/* Join PCG group
 * TODO UPDATE TO GROUP DIGEST WHEN GROUP MODIFIED */
func (p *Peer) JoinGroup(g Group) {
	//fmt.Println(g.String())
	hsha2 := sha256.Sum256(g.Data[:])
	g.AddParticipant(p.node.SocketAddr())
	p.storage[hsha2] = &g
}

func (p *Peer) String() string {
	str := ""
	for _, g := range p.Groups() {
		str = str + g.String()
	}
	return str
}
