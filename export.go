package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type JsonNode struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type JsonEdge struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type JsonOutput struct {
	Nodes []JsonNode `json:"nodes"`
	Edges []JsonEdge `json:"edges"`
}

func Jsonify(ng *NetworkGraph) []byte {
	var out_nodes []JsonNode
	var out_edges []JsonEdge

	for _, node := range ng.GetNodes() {
		out_nodes = append(out_nodes, JsonNode{
			Id:    string(node.GetIP()),
			Label: string(node.GetIP()),
		})

		fmt.Sprintf("%s", node)

	}

	for index, edge := range ng.GetEdges() {
		out_edges = append(out_edges, JsonEdge{
			Id:     strconv.Itoa(index),
			Source: string(edge.LocalAddress),
			Target: string(edge.RemoteAddress),
		})
	}

	output := JsonOutput{out_nodes, out_edges}

	b, _ := json.MarshalIndent(output, "", "  ")
	return b
}
