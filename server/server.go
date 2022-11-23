package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	auction "github.com/adamjhr/ds2-handin-5/proto"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The port of this process")
)

type Server struct {
	auction.FrontendToServerServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Listenered Initialised")

	grpcServer := grpc.NewServer()
	log.Println("Server declared")

	auction.RegisterFrontendToServerServer(grpcServer, &Server{})
	log.Println("Server registered")

	grpcServer.Serve(lis)
	log.Println("Server initialised")
}

// ------------ MESSAGES ------------ \\

func (c *Server) FrontendBid(ctx context.Context, in *auction.FrontendBidRequest) (*auction.FrontendAck, error) {
	log.Println("Bid Request Recieved")
	log.Println(in.Id)
	log.Println(in.Count)
	log.Println(in.Amount)
	return nil, nil
}
