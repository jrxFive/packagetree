package graph

import (
	"errors"
	"fmt"
)

const (
	NODE_NOT_FOUND string = "Node not found"
)

type Node struct {
	name  string
	edges []*Node
}

type TreeGraph struct {
	tree map[string]*Node
}

//Creates an empty TreeGraph type
func NewGraph() (*TreeGraph, error) {

	return &TreeGraph{
		tree: make(map[string]*Node),
	}, nil
}

//If the name exists return the Node, other wise return empty Node struct with error
func (g *TreeGraph) Get(name string) (Node, error) {
	if g.Exists(name) {
		return *g.tree[name], nil
	} else {
		return Node{}, errors.New(fmt.Sprintf("%s:%s", NODE_NOT_FOUND, name))
	}
}

//Return true or false if name exists in graph
func (g *TreeGraph) Exists(name string) bool {
	_, exists := g.tree[name]
	return exists
}

//Performs actual insert/update against graph
func (g *TreeGraph) addNode(node *Node) {
	g.tree[node.name] = node
	g.addEdge(node)
}

//Attempts to add new name, will check that dependencies exist if all exist node will be added,
//otherwise and error is returned
func (g *TreeGraph) Add(name string, edges ...string) error {

	var edgeNodes []*Node

	if len(edges) >= 1 {
		for _, edge := range edges {
			n, err := g.Get(edge)
			if err == nil {
				edgeNodes = append(edgeNodes, &n)
			} else {
				return errors.New(fmt.Sprintf("%s:%s", NODE_NOT_FOUND, edge))
			}
		}
	}

	tempNode := &Node{
		name:  name,
		edges: edgeNodes,
	}

	g.addNode(tempNode)
	return nil
}

//Performs removal of node
func (g *TreeGraph) removeNode(name string) {
	delete(g.tree, name)
}

//Attempts to remove name from graph as long as no dependencies require it
func (g *TreeGraph) Remove(name string) error {

	if g.Exists(name) {
		for _, edge := range g.tree[name].edges {
			if g.Exists(edge.name) {
				for _, e := range edge.edges {
					if e.name == name {
						return errors.New(fmt.Sprintf("Dependency with:%s exists cannot remove:%s", edge.name, name))
					}
				}

			} else {
				continue
			}
		}
		g.removeNode(name)
		return nil
	} else {
		return nil
	}

}

//Links nodes together as edges
func (g *TreeGraph) addEdge(node *Node) {
	var temp []*Node

	for _, edgeNode := range node.edges {
		edgeNode.edges = temp
		edgeNode.edges = append(edgeNode.edges, node)
	}
}
