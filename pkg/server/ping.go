package server

import (
	"context"

	"github.com/scottshotgg/graph-grpc-test/pkg/grapherino"
)

func (s *GraphServer) Ping(ctx context.Context, req *grapherino.PingReq) (*grapherino.PingRes, error) {
	// Client should send their ID
	// You will send back your ID
	return &grapherino.PingRes{
		Id: s.id,
	}, nil
}
