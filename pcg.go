package pcg

import (
	"errors"
	"fmt"
	"github.com/a-shine/butter/node"
	uuid "github.com/nu7hatch/gouuid"
)

// TODO: add the notion of maximum memory for storage

type OverlayPCG struct {
	node    *node.Node
	storage map[uuid.UUID]Group
}

// Block from the node's storage by its UUID. If the block is not found, an empty block with an error is returned.
func (o *OverlayPCG) Group(id string) (Group, error) {
	parsedId, err := uuid.ParseHex(id)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return Group{}, err
	}
	if val, ok := o.storage[*parsedId]; ok {
		return val, nil
	}
	return Group{}, errors.New("block not found")
}

func (o *OverlayPCG) Groups() map[uuid.UUID]Group {
	return o.storage
}

// AddBlock to the node's storage. A UUID is generated for every bit of information added to the network (no update
// functionality yet!). Returns the UUID of the new block as a string.
func (o *OverlayPCG) AddGroup(input string) string {
	var data [4096]byte
	copy(data[:], input)
	// TODO: add the logic to break down the data into blocks if it exceeds the block size
	id, _ := uuid.NewV4()
	o.storage[*id] = Group{
		part:  1,
		parts: 1,
		data:  data,
	}
	return id.String()
}

func NewPCGOverlay(node *node.Node) OverlayPCG {
	return OverlayPCG{
		node:    node,
		storage: make(map[uuid.UUID]Group),
	}
}

// maximise participant diversity by checking ip against ip2location db - try to have participants spread across regions

// three states of group
// - want to find people to joingroup
// - goldilocks group
// - intersect uuid (two groups that have seperated and now know about each other again)
