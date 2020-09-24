package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/scottshotgg/graph-grpc-test/pkg/dijkstra"
)

func main() {
	fmt.Println("hi")

	var err = method1()
	if err != nil {
		log.Fatalln("err:", err)
	}
}

func netMap1() *dijkstra.Graph {
	var graph = dijkstra.NewGraph()

	//Add the 5 verticies
	graph.AddVertex("0")
	graph.AddVertex("1")
	graph.AddVertex("2")
	graph.AddVertex("3")
	graph.AddVertex("4")
	graph.AddVertex("5")

	// Add the arcs
	graph.AddArc("0", "1", 1)
	graph.AddArc("0", "2", 1)
	graph.AddArc("1", "3", 1)
	graph.AddArc("3", "4", 1)
	graph.AddArc("2", "4", 1)
	graph.AddArc("4", "5", 1)

	return graph
}

func netMap2() *dijkstra.Graph {
	var graph = dijkstra.NewGraph()

	// Add the 3 verticies
	graph.AddVertex("0")
	graph.AddVertex("6")
	graph.AddVertex("5")

	// Add the arcs
	graph.AddArc("0", "6", 1)
	graph.AddArc("6", "5", 1)

	return graph
}

func buildNetMap(counts map[string]dijkstra.BestPath) *dijkstra.Graph {
	var graph = dijkstra.NewGraph()

	graph.AddVertex("0")

	for k := range counts {
		graph.AddVertex(k)
	}

	for _, v := range counts {
		var i int

		for {
			var (
				src = v.Path[i]
				dst = v.Path[i+1]
			)

			graph.AddVertex(src)
			graph.AddVertex(dst)

			var err = graph.AddArc(src, dst, 1)
			if err != nil {
				log.Fatalln("err:", err, src, dst)
			}

			err = graph.AddArc(dst, src, 1)
			if err != nil {
				log.Fatalln("err:", err, dst, src)
			}

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

func method1() error {
	fmt.Println("hi")

	var (
		nm1 = netMap1()
		nm2 = netMap2()

		counts = map[string]dijkstra.BestPath{}

		contents, err = nm1.ToGraphViz()
	)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("nm1.gv", []byte(contents), 0777)
	if err != nil {
		return err
	}

	contents, err = nm2.ToGraphViz()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("nm2.gv", []byte(contents), 0777)
	if err != nil {
		return err
	}

	err = calcCounts("1", counts, nm1)
	if err != nil {
		return err
	}

	err = calcCounts("2", counts, nm2)
	if err != nil {
		return err
	}

	var nm3 = buildNetMap(counts)

	err = calcCounts("3", counts, nm3)
	if err != nil {
		return err
	}

	fmt.Println("\nBest:", counts)

	fmt.Println("\ng:", nm3)
	fmt.Println()

	contents, err = nm3.ToGraphViz()
	if err != nil {
		log.Fatalln("err:", err)
	}

	return ioutil.WriteFile("nm3.gv", []byte(contents), 0777)
}

func calcCounts(name string, counts map[string]dijkstra.BestPath, netmap *dijkstra.Graph) error {
	fmt.Println("\nNetMap " + name)

	for i := 1; i < 7; i++ {
		var best, err = netmap.Shortest("0", strconv.Itoa(i))
		if err != nil {
			if err != dijkstra.ErrNoPath {
				log.Fatalln("err:", err)
			}

			continue
		}

		var v, ok = counts[strconv.Itoa(i)]
		if !ok || best.Distance < v.Distance {
			counts[strconv.Itoa(i)] = best
		}

		fmt.Printf("0 -> %d: Shortest distance of %d following path %v\n", i, best.Distance, best.Path)
	}

	return nil
}
