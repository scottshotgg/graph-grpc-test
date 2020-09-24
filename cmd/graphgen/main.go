package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	var (
		nodeCount = 20
		conns     = [][]int{}
		r         = rand.New(rand.NewSource(time.Now().Unix()))
		f, err    = os.Create(strconv.Itoa(nodeCount) + "_nodes_graphviz.gv")
	)

	if err != nil {
		log.Fatalln("err:", err)
	}

	// Init the connMap
	for i := 0; i < nodeCount; i++ {
		conns = append(conns, make([]int, nodeCount))
	}

	for i := 0; i < nodeCount; i++ {
		// Get the connection count
		var connCount = r.Intn(nodeCount)

		j := 0
		for j < connCount {
			var node = r.Intn(nodeCount)

			// Don't allow connections to self
			if node == i {
				continue
			}

			// Check if we already have a connection for that
			if conns[i][node] == 0 {
				fmt.Fprintf(f, "node%d -> node%d [dir=both]\n", i, node)

				conns[i][node] = 1
				conns[node][i] = 1
			}

			j++
		}
	}
}
