package server

import (
	"context"
	"errors"
	"net"

	"github.com/google/uuid"
	"github.com/scottshotgg/graph-grpc-test/dijkstra"
	"github.com/scottshotgg/graph-grpc-test/pkg/grapherino"
	"google.golang.org/grpc"
)

type (
	GraphServer struct {
		netMap *dijkstra.Graph
	}
)

func (g *GraphServer) Exchange(ctx context.Context, req *grapherino.ExchangeReq) (*grapherino.ExchangeRes, error) {
	return nil, errors.New("not implemented")
}

func Start(addr string) error {
	var lis, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	var g = dijkstra.NewGraph()
	g.AddVertex(uuid.New().String())

	var s = GraphServer{
		netMap: g,
	}

	var grpcServer = grpc.NewServer()

	grapherino.RegisterGrapherinoServer(grpcServer, &s)

	return grpcServer.Serve(lis)
}
