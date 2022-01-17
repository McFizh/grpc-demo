package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/mcfizh/grpc-demo/protos"
	"google.golang.org/grpc"
)

type Server struct {
	protos.UnimplementedRandomizerServer
}

var connCnt = 0

func (s *Server) GetValueStream(req *protos.RangeQuery, stream protos.Randomizer_GetValueStreamServer) error {
	timer := time.NewTicker(2 * time.Second)

	connCnt++

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Stream end\n")
			connCnt--
			return nil
		case <-timer.C:
			val := rand.Int31n(req.MaxRange-req.MinRange) + req.MinRange
			log.Printf("(%d) Sending value: %d\n", connCnt, val)
			res := &protos.ValueReply{Value: val}
			stream.Send(res)
		}
	}

}

func (s *Server) GetValuePair(ctx context.Context, req *protos.RangeQuery) (*protos.ValuePairReply, error) {
	var l1 int32
	var numbers = make([]*protos.ValueReply, req.Count)

	for l1 = 0; l1 < req.Count; l1++ {
		val := rand.Int31n(req.MaxRange-req.MinRange) + req.MinRange
		numbers[l1] = &protos.ValueReply{Value: val}
	}

	return &protos.ValuePairReply{
		Values: numbers,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
		return
	}

	grpcServer := grpc.NewServer()
	protos.RegisterRandomizerServer(grpcServer, &Server{})

	log.Printf("Server listening at port 9000\n")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
