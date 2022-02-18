package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/cs347-cw/pcg"
)

func add(overlay *pcg.Peer) {
	fmt.Println("Input information:")
	in := bufio.NewReader(os.Stdin)
	data, _ := in.ReadString('\n') // Read string up to newline
	uuid := pcg.PCGStore(overlay, data)
	clear()
	fmt.Println("UUID:", uuid)
	fmt.Println("Data:", data)
	fmt.Println("Enter to continue...")
	in.ReadString('\n')
	clear()
}

func retrieve(overlay *pcg.Peer) {
	fmt.Println("Information UUID:")
	in := bufio.NewReader(os.Stdin)
	uuid, _ := in.ReadString('\n') // Read string up to newline
	data := pcg.NaiveRetrieve(overlay, uuid)
	clear()
	fmt.Println(string(data))
	fmt.Println("Enter to continue...")
	in.ReadString('\n')
	clear()
}

func printAll(peer *pcg.Peer) {
	fmt.Println(peer.String())
	fmt.Println("Enter to continue...")
	in := bufio.NewReader(os.Stdin)
	in.ReadString('\n')
	clear()
}

func interact(overlayInterface node.Overlay) {
	peer := overlayInterface.(*pcg.Peer)
	fmt.Println("Sock addr: ", peer.Node().SocketAddr())
	for {
		// prompt to pcgStore or pcgRetrieve information
		var interactionType string
		fmt.Print("add(1) or pcgRetrieve(2) or All My IDs(3) information?")
		fmt.Scanln(&interactionType)
		clear()
		switch interactionType {
		case "1":
			add(peer)
		case "2":
			retrieve(peer)
		case "3":
			printAll(peer)
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

func clear() {
	fmt.Print("\033[H\033[2J")
}
