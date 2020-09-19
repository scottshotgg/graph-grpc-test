package main

import (
	"flag"
	"log"
	"strings"

	"github.com/scottshotgg/graph-grpc-test/pkg/server"
)

func main() {
	var (
		peers   = flag.String("peers", "", "")
		addr    = flag.String("addr", "", "")
		errChan = make(chan error)
		n       = server.New()
	)

	flag.Parse()

	if *addr == "" {
		log.Fatalln("Must give an addr")
	}

	log.Printf("Address: \"%s\"\n", *addr)

	// Start second node with first node's addr as an initial peer
	go func() {
		if *peers == "" {
			errChan <- n.Start(*addr)
		} else {
			errChan <- n.Start(*addr, strings.Split(*peers, ",")...)
		}
	}()

	for err := range errChan {
		log.Fatalln("err:", err)
	}
}
