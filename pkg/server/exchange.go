package server

import (
	"context"
	"errors"
	"log"

	"github.com/scottshotgg/graph-grpc-test/pkg/grapherino"
	"google.golang.org/grpc/peer"
)

const stubWeight = 1

func (s *GraphServer) Exchange(ctx context.Context, req *grapherino.ExchangeReq) (*grapherino.ExchangeRes, error) {
	var p, ok = peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("no address available; rejecting request")
	}

	log.Printf("Node \"%s\" at \"%s\" is initiating exchange\n", req.GetId(), p.Addr.String())

	var connections []*grapherino.Connection

	// Fetch yourself from the network map
	v, ok := s.netMap.Verticies[s.id]
	if !ok {
		// Catastrophic error
		// TODO:
		return nil, errors.New("Could not find self")
	}

	// Range over all of your direct connections and build the array
	for to, weight := range v.Arcs() {
		connections = append(connections, &grapherino.Connection{
			To:     to,
			Weight: weight,
		})
	}

	log.Printf("Adding new connection to: \"%s\"\n", req.GetId())

	// TODO: need to check here that the node isn't already in your list
	// Add a new vertex for the node that is connecting
	s.netMap.AddVertex(req.GetId())

	// TODO: need to set the weight somehow
	// Add a new arc for the node that is connection
	var err = s.netMap.AddArc(s.id, req.GetId(), stubWeight)
	if err != nil {
		return nil, err
	}

	// for _, v := range s.netMap.Verticies {
	// 	var connections []*grapherino.Connection
	// nodes[v.ID] = &grapherino.Connections{
	// 	Connections: connections,
	// }
	// }

	// NetMap: &grapherino.NetworkMap{
	// 	Nodes: nodes,
	// },

	s.PrintNetMap()

	return &grapherino.ExchangeRes{
		Id:          s.id,
		Connections: connections,
	}, nil
}
