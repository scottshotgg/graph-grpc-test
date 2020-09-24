package dijkstra

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func FromGraphViz(file string) (*Graph, error) {
	var contents, err = ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var lines = strings.Split(string(contents), "\n")
	fmt.Println("lines", lines)

	var g = NewGraph()

	// Ignore the first and last line
	for _, line := range lines[1 : len(lines)-2] {
		var parts = strings.Split(line, " ")

		if len(parts) < 4 {
			return nil, errors.New("invalid amount of parts")
		}

		var (
			src = parts[0]
			dst = parts[2]
		)

		g.AddVertex(src)
		g.AddVertex(dst)
		g.AddArc(src, dst, 1)
		g.AddArc(src, dst, 1)
	}

	// TODO: fix this

	return nil, errors.New("not implemented")
}

func (g *Graph) ToGraphViz() (string, error) {
	var (
		wrapper = "digraph G {\n%s}"
		seen    = map[string]map[string]int64{}
		buf     = bytes.NewBuffer(nil)
	)

	for k := range g.Verticies {
		seen[k] = map[string]int64{}
	}

	for k, v := range g.Verticies {
		for dst, weight := range v.arcs {
			var _, ok = seen[k][dst]
			if ok {
				continue
			}

			fmt.Fprintf(buf, "\t%s -> %s [dir=both]\n", k, dst)

			seen[k][dst] = weight
			seen[dst][k] = weight
		}
	}

	return fmt.Sprintf(wrapper, buf.String()), nil
}
