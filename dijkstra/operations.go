package dijkstra

import "log"

func (g *Graph) PrintConnections() {
	for _, v := range g.Verticies {
		v.PrintArcs()
	}
}

func (v *Vertex) PrintArcs() {
	for dst, weight := range v.arcs {
		log.Printf("\"%s\" <-> \"%s\" with weight of %d\n", v.ID, dst, weight)
	}
}

func (v *Vertex) Arcs() map[string]int64 {
	return v.arcs
}
