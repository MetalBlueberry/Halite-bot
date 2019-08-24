package astar

import "container/heap"

// astar is an A* pathfinding implementation.

// Pather is an interface which allows A* searching on arbitrary objects which
// can represent a weighted graph.
type Pather interface {
	// PathNeighbors returns the direct neighboring nodes of this node which
	// can be pathed to.
	PathNeighbors() []Pather
	// PathNeighborCost calculates the exact movement cost to neighbor nodes.
	PathNeighborCost(to Pather) float64
	// PathEstimatedCost is a heuristic method for estimating movement costs
	// between non-adjacent nodes.
	PathEstimatedCost(to Pather) float64
}

// node is a wrapper to store A* data for a Pather node.
type node struct {
	pather Pather
	cost   float64
	rank   float64
	parent *node
	open   bool
	closed bool
	index  int
}

// nodeMap is a collection of nodes keyed by Pather nodes for quick reference.
type nodeMap map[Pather]*node

// get gets the Pather object wrapped in a node, instantiating if required.
func (nm nodeMap) get(p Pather) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{
			pather: p,
		}
		nm[p] = n
	}
	return n
}

func reverse(p []Pather) {
	for left, right := 0, len(p)-1; left < right; left, right = left+1, right-1 {
		p[left], p[right] = p[right], p[left]
	}
}

func unwindPath(curr *node) []Pather {
	p := []Pather{}
	for curr != nil {
		p = append(p, curr.pather)
		curr = curr.parent
	}
	return p
}

// Path calculates a short path and the distance between the two Pather nodes.
//
// If no path is found, found will be false.
func Path(from, to Pather, iterations int) (path []Pather, distance float64, found bool, bestPath []Pather) {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	fromNode := nm.get(from)
	fromNode.rank = from.PathEstimatedCost(to)
	fromNode.open = true
	heap.Push(nq, fromNode)
	bestNode := fromNode

	for {
		if nq.Len() == 0 {
			// There's no path, return found false.
			return
		}

		current := heap.Pop(nq).(*node)
		current.open = false
		current.closed = true

		if current == nm.get(to) || iterations == 0 {
			// Found a path to the goal or run out of iterations
			p := unwindPath(current)
			bestPath := unwindPath(bestNode)
			reverse(p)
			reverse(bestPath)
			return p, current.cost, iterations != 0, bestPath
		}
		iterations--

		for _, neighbor := range current.pather.PathNeighbors() {
			cost := current.cost + current.pather.PathNeighborCost(neighbor)
			neighborNode := nm.get(neighbor)
			if cost < neighborNode.cost {
				if neighborNode.open {
					heap.Remove(nq, neighborNode.index)
				}
				neighborNode.open = false
				neighborNode.closed = false
			}
			if !neighborNode.open && !neighborNode.closed {
				neighborNode.cost = cost
				neighborNode.open = true
				neighborNode.rank = cost + neighbor.PathEstimatedCost(to)
				neighborNode.parent = current
				if bestNode.rank-bestNode.cost >= neighborNode.rank-neighborNode.cost {
					bestNode = neighborNode
				}
				heap.Push(nq, neighborNode)
			}
		}
	}
}
