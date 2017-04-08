package main

import (
	"fmt"
	// "sync"
)

type GraphError struct {
	action, message string
}

func (e *GraphError) Error() string {
	return fmt.Sprintf("Error in %s: %s", e.action, e.message)
}

// NetworkGraph

type NetworkGraph struct {
	nodes map[IPAddr]NetworkDevice
	edges []Port
}

func NewNetworkGraph() NetworkGraph {
	return NetworkGraph{
		nodes: make(map[IPAddr]NetworkDevice),
	}
}

func (ng *NetworkGraph) GetNode(id IPAddr) (NetworkDevice, error) {
	if node, ok := ng.nodes[id]; ok {
		return node, nil
	} else {
		return nil, &GraphError{
			"GetNode",
			fmt.Sprintf("Node %s doesn't exist!", id),
		}
	}
}

func (ng *NetworkGraph) GetNodes() map[IPAddr]NetworkDevice {
	return ng.nodes
}

func (ng *NetworkGraph) GetNodesCount() int {
	return len(ng.nodes)
}

func (ng *NetworkGraph) AddNode(addr IPAddr, node NetworkDevice) error {
	if _, ok := ng.nodes[addr]; ok {
		return &GraphError{
			"AddNode",
			fmt.Sprintf("Node with address %s already exists in the graph.", addr),
		}
	} else {
		ng.nodes[addr] = node
		return nil
	}
}

func (ng *NetworkGraph) CrawlOnce(startID IPAddr) []IPAddr {
	startNode, _ := ng.GetNode(startID)

	neighbors := []IPAddr{}

	for port, ip := range startNode.GetNeighbors() {
		if _, node_exists := ng.nodes[ip]; !node_exists {
			ng.AddNode(ip, NewNetworkSwitch(ip))
		}

		ng.edges = append(ng.edges, Port{
			LocalPort:     port,
			RemotePort:    0,
			LocalAddress:  startID,
			RemoteAddress: ip,
		})

		neighbors = append(neighbors, ip)
	}

	return neighbors
}

func (ng *NetworkGraph) GetEdges() []Port {
	return ng.edges
}

func (ng *NetworkGraph) CrawlRecursively(id IPAddr) []IPAddr {
	if _, err := ng.GetNode(id); err != nil {
		return nil
	}

	visited := make(map[IPAddr]bool)
	nodes := []IPAddr{}

	ng.dfsCrawl(id, visited, &nodes)

	return nodes
}

func (ng *NetworkGraph) dfsCrawl(id IPAddr, visited map[IPAddr]bool, nodes *[]IPAddr) {
	if _, ok := visited[id]; ok {
		return
	}

	visited[id] = true
	*nodes = append(*nodes, id)

	startNode, _ := ng.GetNode(id)

	for port, ip := range startNode.GetNeighbors() {
		if _, node_exists := ng.nodes[ip]; !node_exists {
			ng.AddNode(ip, NewNetworkSwitch(ip))
		}

		ng.edges = append(ng.edges, Port{
			LocalPort:     port,
			RemotePort:    0,
			LocalAddress:  id,
			RemoteAddress: ip,
		})

		ng.dfsCrawl(ip, visited, nodes)
	}

}
