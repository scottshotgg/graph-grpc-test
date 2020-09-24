package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/scottshotgg/graph-grpc-test/pkg/server"
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

	var errChan = make(chan error)

	for i, conns := range conns {
		var (
			port = 5000 + i
			addr = ":" + strconv.Itoa(port)
		)

		fmt.Println("addr:", addr)

		var peers []string

		for j, conn := range conns {
			if conn == 0 {
				continue
			}

			peers = append(peers, ":"+strconv.Itoa(5000+j))
		}

		go func(i int, peers []string) {
			fmt.Println("node", i)
			fmt.Println("peers:", peers)

			errChan <- server.New().Start(addr, peers...)
		}(i, peers)

		time.Sleep(2 * time.Second)
	}

	for err := range errChan {
		log.Fatalln("err:", err)
	}
}
