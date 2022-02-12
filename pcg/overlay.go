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

// AddBlock to the node's storage. A UUID is generated for every bit of information added to the network (no update
// functionality yet!). Returns the UUID of the new block as a string.
func (o *PCG) AddBlock(data string) string {
	// TODO: add the logic to break down the data into blocks if it exceeds the block size
	hsha2 := sha256.Sum256([]byte(data))
	var formattedData [4096]byte
	copy(formattedData[:], data)
	o.storage[hsha2] = Group{
		data: formattedData,
	}
	return fmt.Sprintf("%x", hsha2)
}

// Group from the node's storage by its UUID. If the block is not found, an empty block with an error is returned.
func (o *PCG) Block(id string) (Group, error) {
	var hash [32]byte
	data, _ := hex.DecodeString(id)
	copy(hash[:], data)
	if val, ok := o.storage[hash]; ok {
		return val, nil
	}
	return Group{}, errors.New("block not found")
}

func NewOverlay(node *node.Node) PCG {
	return PCG{
		node:    node,
		storage: make(map[[32]byte]Group),
	}
}
