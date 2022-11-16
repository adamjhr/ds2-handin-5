package main

import (
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
	auction.UnimplementedTemplateServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	auction.RegisterTemplateServer(grpcServer, &Server{})

	grpcServer.Serve(lis)
}
