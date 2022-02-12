package pcg

import (
	"fmt"
	"github.com/a-shine/butter/discover"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/traverse"
	"os"
	"os/signal"
	"syscall"
)

// Spawn node into the network (the node serves as an entry-point to the butter network). You can also do this manually
// to have more control over the specific protocols used in your dapp. This function presents a simple abstraction with
// the included default butter protocols.
func Spawn(overlay node.Overlay, traverseFlag bool) {
	n := overlay.Node()
	setupLeaveHandler(n)
	go discover.Discover(overlay)
	if traverseFlag {
		go traverse.Traverse(n)
	}
	n.Start(overlay)
}

func SpawnOverlay(node *node.Node, traverseFlag bool) {
	overlay := NewOverlay(node) // Creates a new overlay network
	AppendRetrieveBehaviour(overlay.Node())
	Spawn(&overlay, traverseFlag)
}

// setupLeaveHandler creates a listener on a new goroutine which will notify the program if it receives an interrupt
// from the OS and then handles the node leaving the network gracefully.
func setupLeaveHandler(node *node.Node) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\rLeaving the butter network...")
		node.Shutdown()
		os.Exit(0)
	}()
}
