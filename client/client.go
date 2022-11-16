package main

import (
	"flag"
	"log"

	"github.com/adamjhr/ds2-handin-5/auction"

	"google.golang.org/grpc"
)

var (
	port     = flag.Int("port", 9000, "The port of this process")
	replicas = make([]*grpc.ClientConn, len(flag.Args()))
)

func main() {
	flag.Parse()
	replicaPorts := flag.Args()
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	for i, _ := range replicaPorts {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9080", grpc.WithInsecure())
		defer conn.Close()
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		replicas[i] = conn
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)

	//  Create new Client from generated gRPC code from proto
	c := auction.TemplateClient(conn)

	for {

	}
}
