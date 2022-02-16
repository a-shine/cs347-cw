package main

import (
	"bufio"
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/cs347-cw/pcg"
	"os"
)

func add(overlay *pcg.PCG) {
	fmt.Println("Input information:")
	in := bufio.NewReader(os.Stdin)
	data, _ := in.ReadString('\n') // Read string up to newline
	uuid := pcg.PCGStore(overlay, data)
	fmt.Println("UUID:", uuid)
}

func retrieve(overlay *pcg.PCG) {
	fmt.Println("Information UUID:")
	in := bufio.NewReader(os.Stdin)
	uuid, _ := in.ReadString('\n') // Read string up to newline
	data := pcg.NaiveRetrieve(overlay, uuid)
	fmt.Println(string(data))
}

func interact(overlayInterface node.Overlay) {
	pcg := overlayInterface.(*pcg.PCG)
	for {
		// prompt to pcgStore or pcgRetrieve information
		var interactionType string
		fmt.Print("add(1) or pcgRetrieve(2) information?")
		fmt.Scanln(&interactionType)

		switch interactionType {
		case "1":
			add(pcg)
		case "2":
			retrieve(pcg)
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func main() {
	butterNode, _ := node.NewNode(0, 512, false)
	butterNode.RegisterClientBehaviour(interact)

	overlay := pcg.NewPCG(&butterNode, 512) // Creates a new overlay network
	pcg.AppendRetrieveBehaviour(overlay.Node())
	pcg.AppendGroupStoreBehaviour(overlay.Node())

	butter.Spawn(&overlay, false)
}
