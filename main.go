package main

import (
	"fmt"

	"github.com/scottshotgg/graph-grpc-test/dijkstra"
)

func main() {
	fmt.Println("hi")

	method1()
}

func netMap1() *dijkstra.Graph {
	var graph = dijkstra.NewGraph()

	//Add the 3 verticies
	graph.AddVertex(0)
	graph.AddVertex(1)
	graph.AddVertex(2)
	graph.AddVertex(3)
	graph.AddVertex(4)
	graph.AddVertex(5)

	//Add the arcs
	graph.AddArc(0, 1, 1)
	graph.AddArc(0, 2, 1)
	graph.AddArc(1, 3, 1)
	graph.AddArc(3, 4, 1)
	graph.AddArc(2, 4, 1)
	graph.AddArc(4, 5, 1)

	return graph
}

func netMap2() *dijkstra.Graph {
	var graph = dijkstra.NewGraph()

	//Add the 3 verticies
	graph.AddVertex(0)
	graph.AddVertex(6)
	graph.AddVertex(5)

	//Add the arcs
	graph.AddArc(0, 6, 1)
	graph.AddArc(6, 5, 1)

	return graph
}

func buildNetMap(nm map[int]dijkstra.BestPath) *dijkstra.Graph {
	var graph = dijkstra.NewGraph()

	for k, v := range nm {
		graph.AddVertex(k)

		var i int

		for {
			graph.AddArc(v.Path[i], v.Path[i+1], 1)
			i++

			if i+1 < len(v.Path) {
				continue
			}

			break
		}
	}

	return graph
}

// TODO: might need this later; find the largest ID node in the graph
// func findLargest()

func method1() {
	fmt.Println("hi")

	var (
		nm1 = netMap1()
		nm2 = netMap2()

		counts = map[int]dijkstra.BestPath{}
	)

	fmt.Println("NetMap1")

	for i := 1; i < 7; i++ {
		best, err := nm1.Shortest(0, i)
		if err != nil {
			continue
		}

		counts[i] = best

		fmt.Printf("0 -> %d: Shortest distance of %d following path %v\n", i, best.Distance, best.Path)
	}

	fmt.Println()
	fmt.Println("NetMap2")

	for i := 1; i < 7; i++ {
		best, err := nm2.Shortest(0, i)
		if err != nil {
			continue
		}

		var _, ok = counts[i]
		if !ok || best.Distance < counts[i].Distance {
			counts[i] = best
		}

		fmt.Printf("0 -> %d: Shortest distance of %d following path %v\n", i, best.Distance, best.Path)
	}

	fmt.Println("Best:", counts)

	var nm3 = buildNetMap(counts)

	fmt.Println()
	fmt.Println("NetMap3")

	for i := 1; i < 7; i++ {
		best, err := nm3.Shortest(0, i)
		if err != nil {
			continue
		}

		var _, ok = counts[i]
		if !ok || best.Distance < counts[i].Distance {
			counts[i] = best
		}

		fmt.Printf("0 -> %d: Shortest distance of %d following path %v\n", i, best.Distance, best.Path)
	}
}
