package pcg

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/a-shine/butter/node"
)

type PCG struct {
	node    *node.Node
	storage map[[32]byte]Group
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

func NewPCG(node *node.Node) PCG {
	return PCG{
		node:    node,
		storage: make(map[[32]byte]Group),
	}
}
