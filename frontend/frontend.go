package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	auction "github.com/adamjhr/ds2-handin-5/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port           = flag.Int("port", 7000, "The port of this process")
	replicas       = make([]*grpc.ClientConn, len(flag.Args()))
	replicaClients = make([]auction.FrontendToServerClient, 1)
	count          = 0
)

func main() {
	flag.Parse()
	replicaPorts := flag.Args()
	// Create a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	for _, p := range replicaPorts {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(fmt.Sprintf(":%v", p), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		replicas = append(replicas, conn)
		replicaClients = append(replicaClients, auction.NewFrontendToServerClient(conn))
	}
	log.Println("About to send mulitcast request")
	MultiCastBid(10)
	log.Println("Multicast request sent")
	for {

	}
}

func MultiCastBid(amount int) {
	for i := range replicaClients {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		replicaClients[i].FrontendBid(ctx, &auction.FrontendBidRequest{Id: int32(*port), Count: int32(count), Amount: int32(amount)})
	}
}

func MultiCastResult() {

}
