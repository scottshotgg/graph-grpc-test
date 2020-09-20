package server

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/scottshotgg/graph-grpc-test/pkg/dijkstra"
	"github.com/scottshotgg/graph-grpc-test/pkg/grapherino"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const stubWeight = 1

func (s *GraphServer) Exchange(ctx context.Context, req *grapherino.ExchangeReq) (*grapherino.ExchangeRes, error) {
	var p, ok = peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("no address available; rejecting request")
	}

	var addr = req.GetAddr()

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

	// Make a new client conn
	var conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// Make a new client to the connection
	var c = grapherino.NewGrapherinoClient(conn)

	log.Printf("Pinging peer: %s\n\n", addr)

	var start = time.Now().UnixNano()

	// Ping the client that is trying to connect
	res, err := c.Ping(ctx, &grapherino.PingReq{
		Id: s.id,
	})

	var elapsed = time.Now().UnixNano() - start

	if err != nil {
		log.Println("Could not establish connection to peer:", addr)
		// TODO:
		return nil, err
	}

	log.Printf("Successfully pinged peer: \"%s\" (\"%s\") with latency of %d\n", addr, res.GetId(), elapsed)

	// If you can do the exchange, add them as a peer
	s.peers[res.GetId()] = &Peer{
		id:     res.GetId(),
		addr:   addr,
		conn:   conn,
		client: c,
		err:    err,
	}

	// TODO: need to check here that the node isn't already in your list
	// Add a new vertex for the node that is connecting
	s.netMap.AddVertex(req.GetId())

	// TODO: need to set the weight somehow
	// Add a new arc for the node that is connection
	err = s.netMap.AddArc(s.id, req.GetId(), elapsed)
	if err != nil {
		return nil, err
	}

	var counts = map[string]dijkstra.BestPath{}

	// Calculate the shortest path to every node with the new connections that we got back
	s.calcCounts(counts)

	// Build the new network map from the shortest path graph
	s.buildNetMap(counts)

	// Print the network map
	s.PrintNetMap()

	return &grapherino.ExchangeRes{
		Id:          s.id,
		Connections: connections,
	}, nil
}

func (s *GraphServer) determineWeight(ctx context.Context, typeOf, id string) (int64, error) {
	switch typeOf {
	case "ping":
		var peer, ok = s.peers[id]
		if !ok {
			return 0, errors.New("could not find client for peer")
		}

		var (
			start  = time.Now().UnixNano()
			_, err = peer.client.Ping(ctx, &grapherino.PingReq{
				Id: s.id,
			})
		)

		if err != nil {
			return 0, err
		}

		return time.Now().UnixNano() - start, nil

	default:
		return stubWeight, nil
	}
}
