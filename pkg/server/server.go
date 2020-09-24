package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/scottshotgg/graph-grpc-test/pkg/dijkstra"
	"github.com/scottshotgg/graph-grpc-test/pkg/grapherino"
	"google.golang.org/grpc"
)

type (
	GraphServer struct {
		sync.Mutex

		id     string
		netMap *dijkstra.Graph
		peers  map[string]*Peer
	}
)

func New() (s *GraphServer) {
	var (
		netMap = dijkstra.NewGraph()
		id     = uuid.New().String()
	)

	log.Printf("ID: \"%s\"\n", id)

	netMap.AddVertex(id)

	return &GraphServer{
		id:     id,
		netMap: netMap,
		peers:  map[string]*Peer{},
	}
}

type (
	Peer struct {
		id   string
		addr string

		conn   *grpc.ClientConn
		client grapherino.GrapherinoClient
		err    error
	}
)

func (s *GraphServer) Start(addr string, initialPeers ...string) error {
	var lis, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	var grpcServer = grpc.NewServer()

	grapherino.RegisterGrapherinoServer(grpcServer, s)

	go func() {
		for {
			err = s.mapPeers(initialPeers)
			if err != nil {
				return
			}

			time.Sleep(5 * time.Second)
		}
	}()

	fmt.Println("addr started:", addr)

	return grpcServer.Serve(lis)
}

func (s *GraphServer) mapPeers(initPeers []string) error {
	if len(initPeers) != 0 {
		log.Println("Beginning network map replication from peers ...")
	}

	for _, addr := range initPeers {
		log.Println("Dialing peer:", addr)

		// Make a new client conn
		var conn, err = grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			return err
		}

		// Make a new client to the connection
		var c = grapherino.NewGrapherinoClient(conn)

		log.Printf("Exchanging with peer: %s\n\n", addr)

		var ctx = context.Background()

		// Exchange network maps
		res, err := c.Exchange(ctx, &grapherino.ExchangeReq{
			Id:   s.id,
			Addr: addr,
		})

		if err != nil {
			log.Println("Could not establish connection to peer:", addr)
			// TODO:
			continue
		}

		log.Printf("Successfully replicated from peer: \"%s\" (\"%s\")\n", addr, res.GetId())

		// If you can do the exchange, add them as a peer
		s.peers[res.GetId()] = &Peer{
			id:     res.GetId(),
			addr:   addr,
			conn:   conn,
			client: c,
			err:    err,
		}

		// var peer, ok = s.netMap.Verticies[res.GetId()]

		// // If it is already in the map then we have already seen them
		// // Will probably have to handle this one and check things because we
		// // shouldn't get them twice - atleast not in this step
		// if ok {
		// 	// TODO:
		// 	log.Fatalln("handle this double boi shit")
		// }

		// Add a vertex for the new peer (node)
		s.netMap.AddVertex(res.GetId())

		weight, err := s.determineWeight(ctx, "ping", res.GetId())
		if err != nil {
			return err
		}

		// fmt.Println("weight:", weight)

		// Add an arc for the new peer (node)
		// TODO: get the weight from a call timer or benchmark or something
		err = s.netMap.AddArc(s.id, res.GetId(), weight)
		if err != nil {
			return err
		}

		for _, conn := range res.GetConnections() {
			// Add a vertex for the new peer (node)
			s.netMap.AddVertex(conn.GetTo())

			// Add an arc for the new peer (node)
			// TODO: get the weight from a call timer or benchmark or something
			err = s.netMap.AddArc(res.GetId(), conn.GetTo(), weight+conn.GetWeight())
			if err != nil {
				return err
			}
		}
	}

	if len(s.netMap.Verticies) > 0 {
		var counts = map[string]dijkstra.BestPath{}

		// Calculate the shortest path to every node with the new connections that we got back
		s.calcCounts(counts)

		// fmt.Println("counts:", counts)

		// Build the new network map from the shortest path graph
		s.buildNetMap(counts)

		// Print the network map
		s.PrintNetMap()
	}

	return nil
}

func (s *GraphServer) PrintNetMap() {
	log.Println("Connections:")
	s.netMap.PrintConnections()
	log.Println()
}

func (s *GraphServer) calcCounts(counts map[string]dijkstra.BestPath) error {
	for k := range s.netMap.Verticies {
		var best, err = s.netMap.Shortest(s.id, k)
		if err != nil {
			if err != dijkstra.ErrNoPath {
				log.Fatalln("err:", err)
			}

			continue
		}

		var v, ok = counts[k]
		if !ok || best.Distance < v.Distance {
			counts[k] = best
		}

		// log.Printf("%s <-> %s: Shortest distance of %d following path %v\n", s.id, k, best.Distance, best.Path)
	}

	return nil
}

func (s *GraphServer) buildNetMap(counts map[string]dijkstra.BestPath) error {
	var graph = dijkstra.NewGraph()

	graph.AddVertex(s.id)

	for k := range counts {
		graph.AddVertex(k)
	}

	for _, v := range counts {
		var i int

		for {
			var (
				src = v.Path[i]
				dst = v.Path[i+1]
			)

			graph.AddVertex(src)
			graph.AddVertex(dst)

			var err = graph.AddArc(src, dst, v.Distance)
			if err != nil {
				return err
			}

			i++

			if i+1 < len(v.Path) {
				continue
			}

			break
		}
	}

	s.Lock()

	s.netMap = graph

	s.Unlock()

	return nil
}
