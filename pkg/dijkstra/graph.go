package dijkstra

import (
	"errors"

	"github.com/google/uuid"
)

//Graph contains all the graph details
type Graph struct {
	best        int64
	visitedDest bool
	//slice of all verticies available
	Verticies       map[string]*Vertex
	visiting        dijkstraList
	mapping         map[string]string
	usingMap        bool
	highestMapIndex string
}

//NewGraph creates a new empty graph
func NewGraph() *Graph {
	return &Graph{
		Verticies: map[string]*Vertex{},
		mapping:   map[string]string{},
	}
}

//AddNewVertex adds a new vertex at the next available index
func (g *Graph) AddNewVertex() *Vertex {
	for i, v := range g.Verticies {
		if i != v.ID {
			g.Verticies[i] = &Vertex{
				ID: i,
			}

			return g.Verticies[i]
		}
	}

	return g.AddVertex(uuid.New().String())
}

//AddVertex adds a single vertex
func (g *Graph) AddVertex(ID string) *Vertex {
	var v, ok = g.Verticies[ID]
	if ok {
		return v
	}

	g.AddVerticies(Vertex{
		ID: ID,
	})

	return g.Verticies[ID]
}

//GetVertex gets the reference of the specified vertex. An error is thrown if
// there is no vertex with that index/ID.
func (g *Graph) GetVertex(ID string) (*Vertex, error) {
	var v, ok = g.Verticies[ID]
	if !ok {
		return nil, errors.New("Vertex not found")
	}

	return v, nil
}

//SetDefaults sets the distance and best node to that specified
func (g *Graph) setDefaults(Distance int64, BestNode string) {
	for i := range g.Verticies {
		g.Verticies[i].bestVerticies = []string{BestNode}
		g.Verticies[i].distance = Distance
	}
}
