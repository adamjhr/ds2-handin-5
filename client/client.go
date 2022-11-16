package main

import (
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port     = flag.Int("port", 9000, "The port of this process")
	replicas = make([]*grpc.ClientConn, len(flag.Args()))
)

func main() {
	flag.Parse()
	replicaPorts := flag.Args()
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	for _, p := range replicaPorts {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(fmt.Sprintf(":%v", p), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		replicas = append(replicas, conn)
	}

	for {

	}
}
