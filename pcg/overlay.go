package pcg

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/a-shine/butter/node"
)

type PCG struct {
	node       *node.Node
	maxStorage uint64
	storage    map[[32]byte]Group
}

func (o *PCG) Node() *node.Node {
	return o.node
}

// AddGroup to the node's storage. A UUID is generated for every bit of information added to the network (no update
// functionality yet!). Returns the UUID of the new block as a string.
func (o *PCG) AddGroup(data string) string {
	// check if node memory (allocated by the user on initialization) is full - node API
	hsha2 := sha256.Sum256([]byte(data))
	var formattedData [4096]byte
	copy(formattedData[:], data)
	o.storage[hsha2] = NewGroup(formattedData, o.node.SocketAddr())
	return fmt.Sprintf("%x", hsha2)
}

// Group from the node's storage by its UUID. If the block is not found, an empty block with an error is returned.
func (o *PCG) Group(id string) (Group, error) {
	var hash [32]byte
	data, _ := hex.DecodeString(id)
	copy(hash[:], data)
	if val, ok := o.storage[hash]; ok {
		return val, nil
	}
	return Group{}, errors.New("block not found")
}

func (o *PCG) Groups() map[[32]byte]Group {
	return o.storage
}

func MbToBytes(mb uint64) uint64 {
	return uint64(mb * 1024 * 1024)
}

func MaxStorage(maxMemory uint64) uint64 {
	return maxMemory / uint64(GroupSize)
}

func NewPCG(node *node.Node, maxMemoryMb uint64) PCG {
	maxMemory := MbToBytes(maxMemoryMb)
	maxStorage := MaxStorage(maxMemory)
	fmt.Println("Max storage:", maxStorage)
	return PCG{
		node:       node,
		maxStorage: maxStorage,
		storage:    make(map[[32]byte]Group),
	}
}
