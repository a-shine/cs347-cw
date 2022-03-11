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

func (p *Peer) Node() *node.Node { // Required to implement the Butter overlay interface
	return p.node
}

func (p *Peer) AvailableStorage() uint64 { // Required to implement the Butter overlay interface
	return p.maxStorage - p.currentStorage
}

// Group from node's storage by its UUID. If the node is not part of that group, returns nil with an error.
func (p *Peer) Group(id string) (*Group, error) {
	var hash [32]byte
	data, _ := hex.DecodeString(id)
	copy(hash[:], data)
	if group, ok := p.storage[hash]; ok {
		return group, nil
	}
	return nil, errors.New("block not found")
}

// Groups from node's storage by their UUIDs
func (p *Peer) Groups() map[[32]byte]*Group {
	return p.storage
}

// --- Node Group creation and joining ---

// CreateGroup with node. A UUID is generated for every bit of information added to the network (no update
// functionality yet!). Returns the UUID of the new group as a string.
func (p *Peer) CreateGroup(data string) string {
	// TODO: check if node memory (allocated by the user on initialization) is full - node API
	var formattedData [4096]byte
	copy(formattedData[:], data)
	hsha2 := sha256.Sum256(formattedData[:])
	p.storage[hsha2] = NewGroup(formattedData, p.node.SocketAddr())
	p.currentStorage += 4096 //TODO
	return fmt.Sprintf("%x", hsha2)
}

// JoinGroup from a node
func (p *Peer) JoinGroup(g Group) {
	//TODO: UPDATE TO GROUP DIGEST WHEN GROUP MODIFIED
	hsha2 := sha256.Sum256(g.Data[:])
	err := g.AddParticipant(p.node.SocketAddr())
	if err != nil {
		fmt.Println("Unable to join group:", err)
	}
	p.storage[hsha2] = &g
}

// --- Encoders ---

// String of node's groups
func (p *Peer) String() string {
	str := ""
	for _, g := range p.Groups() {
		str = str + g.String()
	}
	return str
}
