package main

import (
	"fmt"
	"log"

	"github.com/scottshotgg/graph-grpc-test/pkg/server"
)

func main() {
	fmt.Println("sup world")

	var err = server.Start(":5001")
	if err != nil {
		log.Println("err:", err)
	}
}
