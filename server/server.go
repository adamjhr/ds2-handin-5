package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	auction "github.com/adamjhr/ds2-handin-5/proto"
	"google.golang.org/grpc"
)

var (
	port              = flag.Int("port", 8000, "The port of this process")
	auctionRequest    = false
	auctionIsFinished = true
	highestBid        int
	highestBidder     int
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
	log.Println("Listener Initialised")

	grpcServer := grpc.NewServer()
	log.Println("Server declared")

	auction.RegisterFrontendToServerServer(grpcServer, &Server{})
	log.Println("Server registered")

	// Serve the replica
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()
	// Listen for when a new auction should be started
	for {
		if auctionRequest && auctionIsFinished {
			auctionRequest = false
			time.Sleep(5 * time.Second)
			NewAuction()
		}
	}
}

// ------------ MESSAGES ------------ \\

func (c *Server) FrontendNewAuction(ctx context.Context, in *auction.FrontendNewAuctionRequest) (*auction.FrontendNewAuctionReply, error) {
	log.Println("New Auction Request Recieved")
	if auctionRequest {
		return &auction.FrontendNewAuctionReply{Id: in.Id, Count: in.Count, Outcome: auction.Outcome_Fail}, nil
	}
	auctionRequest = true
	return &auction.FrontendNewAuctionReply{Id: in.Id, Count: in.Count, Outcome: auction.Outcome_Success}, nil
}

func (c *Server) FrontendBid(ctx context.Context, in *auction.FrontendBidRequest) (*auction.FrontendAck, error) {
	log.Println("Bid Request Recieved")
	// Check the size of the bid and wether an auction is happening
	if in.Amount > int32(highestBid) && !auctionIsFinished {
		highestBid = int(in.Amount)
		log.Println(highestBidder)
		highestBidder = int(in.Id)
		log.Println(highestBidder)
		log.Println(int32(highestBidder))
		return &auction.FrontendAck{Id: in.Id, Count: in.Count, Outcome: auction.Outcome_Success}, nil
	} else {
		return &auction.FrontendAck{Id: in.Id, Count: in.Count, Outcome: auction.Outcome_Fail}, nil
	}
}

func (c *Server) FrontendResult(ctx context.Context, in *auction.FrontendResultRequest) (*auction.FrontendResultReply, error) {
	log.Println("Result Request Recieved")
	return &auction.FrontendResultReply{Id: in.Id, Count: in.Count, Bidder: int32(highestBidder), Amount: int32(highestBid), IsFinished: auctionIsFinished}, nil
}

func NewAuction() {
	log.Println("New Auction Created")
	auctionIsFinished = false
	highestBid = 10
	time.Sleep(time.Second * 30)
	log.Println("Auction has Ended")
	auctionIsFinished = true
}
