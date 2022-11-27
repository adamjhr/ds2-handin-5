package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	auction "github.com/adamjhr/ds2-handin-5/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port           = flag.Int("port", 7000, "The port of this process")
	replicas       = make([]*grpc.ClientConn, len(flag.Args()))
	replicaClients = make([]auction.FrontendToServerClient, 0)
	count          = 0
)

type Server struct {
	auction.ClientToFrontendServer
}

func main() {
	flag.Parse()
	replicaPorts := flag.Args()

	for _, p := range replicaPorts {
		conn, err := grpc.Dial(fmt.Sprintf(":%v", p), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		replicas = append(replicas, conn)
		replicaClients = append(replicaClients, auction.NewFrontendToServerClient(conn))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Listener Initialised")

	grpcServer := grpc.NewServer()
	log.Println("Frontend declared")

	auction.RegisterClientToFrontendServer(grpcServer, &Server{})
	log.Println("Frontend registered")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func (c *Server) ClientNewAuction(ctx context.Context, in *auction.ClientNewAuctionRequest) (*auction.ClientNewAuctionReply, error) {
	count++
	var serverResponse *auction.FrontendNewAuctionReply
	for i := range replicaClients {
		go func(index int) {
			ctx, cancel := context.WithCancel(context.Background())
			var err error
			serverResponse, err = replicaClients[index].FrontendNewAuction(ctx, &auction.FrontendNewAuctionRequest{Id: int32(*port), Count: int32(count)})
			if err != nil {
				log.Printf("Error sending message to server at port: %v", index)
			}
			cancel()
		}(i)
	}
	for {
		if serverResponse != nil {
			log.Println(serverResponse.Id)
			return &auction.ClientNewAuctionReply{Id: serverResponse.Id, Outcome: serverResponse.Outcome}, nil
		}
	}
}

func (c *Server) ClientBid(ctx context.Context, in *auction.ClientBidRequest) (*auction.ClientAck, error) {
	count++
	var serverResponse *auction.FrontendAck
	for i := range replicaClients {
		go func(index int) {
			ctx, cancel := context.WithCancel(context.Background())
			var err error
			serverResponse, err = replicaClients[index].FrontendBid(ctx, &auction.FrontendBidRequest{Id: int32(in.Id), Count: int32(count), Amount: int32(in.Amount)})
			if err != nil {
				log.Printf("Error sending message to server at port: %v", index)
			}
			cancel()
		}(i)
	}
	for {
		if serverResponse != nil {
			log.Println(serverResponse.Id)
			return &auction.ClientAck{Id: serverResponse.Id, Outcome: serverResponse.Outcome}, nil
		}
	}
}

func (c *Server) ClientResult(ctx context.Context, in *auction.ClientResultRequest) (*auction.ClientResultReply, error) {
	count++
	var serverResponse *auction.FrontendResultReply
	for i := range replicaClients {
		go func(index int) {
			ctx, cancel := context.WithCancel(context.Background())
			var err error
			serverResponse, err = replicaClients[index].FrontendResult(ctx, &auction.FrontendResultRequest{Id: in.Id, Count: int32(count)})
			if err != nil {
				log.Printf("Error sending message to server at port: %v", index)
			}
			cancel()
		}(i)
	}
	for {
		if serverResponse != nil {
			log.Println(serverResponse.Id)
			return &auction.ClientResultReply{Id: serverResponse.Id, Amount: serverResponse.Amount, Bidder: serverResponse.Bidder, IsFinished: serverResponse.IsFinished}, nil
		}
	}
}
