package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	fmt.Println("hi")

	var contents, err = ioutil.ReadFile("20_nodes_graphviz.gv")
	if err != nil {
		log.Fatalln("err:", err)
	}

	var lines = strings.Split(string(contents), "\n")
	fmt.Println("lines", lines)

	// Ignore the first and last line
	for _, line := range lines[1 : len(lines)-2] {
		var parts = strings.Split(line, " ")

		if len(parts) < 4 {
			log.Fatalln("invalid amount of parts;", len(parts))
		}

		var (
			src = parts[0]
			dst = parts[2]
		)

		fmt.Printf("src (%s) goes to dst (%s)\n", src, dst)
	}
}
