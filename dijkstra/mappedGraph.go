package dijkstra

import (
	"errors"
	"fmt"
)

//GetMapped gets the key assosciated with the mapped int
func (g *Graph) GetMapped(a string) (string, error) {
	if !g.usingMap || g.mapping == nil {
		return "", errors.New("Map is not being used/initialised")
	}

	for k, v := range g.mapping {
		if v == a {
			return k, nil
		}
	}

	return "", errors.New(fmt.Sprint(a, " not found in mapping"))
}

//GetMapping gets the index associated with the specified key
func (g *Graph) GetMapping(a string) (string, error) {
	// if g.mapping != nil && g.usingMap {
	// 	var b, ok = g.mapping[a]
	// 	if ok {
	// 		return b, nil
	// 	}
	// }

	if !g.usingMap || g.mapping == nil {
		return "", errors.New("Map is not being used/initialised")
	}

	var b, ok = g.mapping[a]
	if ok {
		return b, nil
	}

	return "", errors.New(fmt.Sprint(a, " not found in mapping"))
}

//AddMappedVertex adds a new Vertex with a mapped ID (or returns the index if
// ID already exists).
func (g *Graph) AddMappedVertex(ID string) string {
	if !g.usingMap || g.mapping == nil {
		g.usingMap = true
		g.mapping = map[string]string{}
		g.highestMapIndex = ""
	}

	i, ok := g.mapping[ID]
	if ok {
		return i
	}

	i = g.highestMapIndex
	// g.highestMapIndex++
	g.mapping[ID] = i

	return g.AddVertex(i).ID
}

//AddMappedArc adds a new Arc from Source to Destination, for when verticies are
// referenced by strings.
func (g *Graph) AddMappedArc(Source, Destination string, Distance int64) error {
	return g.AddArc(g.AddMappedVertex(Source), g.AddMappedVertex(Destination), Distance)
}

//AddArc is the default method for adding an arc from a Source Vertex to a
// Destination Vertex
func (g *Graph) AddArc(Source, Destination string, Distance int64) error {
	// if len(g.Verticies) <= Source || len(g.Verticies) <= Destination {
	// 	return errors.New("Source/Destination not found")
	// }

	var v, ok = g.Verticies[Source]
	if !ok {
		return errors.New("Source/Destination not found")
	}

	v.AddArc(Destination, Distance)

	return nil
}
