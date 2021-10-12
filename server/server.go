package main

import (
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

func (s *Server) GetValues(req *protos.RangeQuery, stream protos.Randomizer_GetValuesServer) error {
	timer := time.NewTicker(2 * time.Second)

	connCnt++
	log.Printf("Min: %d, Max: %d\n", req.MinRange, req.MaxRange)

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
