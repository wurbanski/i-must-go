package main

import (
	"fmt"
	// "github.com/gonum/graph"
	// "github.com/gonum/graph/simple"
)

// var ports = make(map[int]*Port)
var switches = make(map[IPAddr]*NetworkSwitch)

func main() {

	startIP := IPAddr("10.1.1.0")
	startSwitch := NewNetworkSwitch(startIP)

	network := NewNetworkGraph()
	network.AddNode(startIP, startSwitch)
	network.CrawlRecursively(startIP)
	fmt.Printf("Found %d nodes!\n", network.GetNodesCount())
	fmt.Printf("%s", network.GetEdges())

}
