package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	auction "github.com/adamjhr/ds2-handin-5/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port     = flag.Int("port", 6999, "The port of this process")
	receiver = flag.Int("receiver", 7000, "The port of the corresponding frontend")
)

func main() {
	flag.Parse()

	// Dial frontend
	conn, err := grpc.Dial(fmt.Sprintf(":%v", *receiver), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	c := auction.NewClientToFrontendClient(conn)

	defer conn.Close()

	fmt.Println(" ------------------------------- ")
	fmt.Println(" --- WELCOME TO THIS AUCTION --- ")
	fmt.Println(" ------------------------------- ")
	fmt.Println(" To make a bid, type 'bid <amount>'")
	fmt.Println(" To see the results of the auction, type 'result'")
	fmt.Println(" To start a new auction, type 'auction'")

	scanner := bufio.NewScanner(os.Stdin)

	// Loop listening for console input
	for {
		if scanner.Scan() {
			line := scanner.Text()
			lineSplit := strings.Split(line, " ")
			if len(lineSplit) != 0 {
				ctx, cancel := context.WithCancel(context.Background())
				if lineSplit[0] == "bid" {
					// convert bid-amount to number
					bidAmount, errNum := strconv.Atoi(lineSplit[1])
					if errNum != nil {
						log.Println("Invalid bid amount")
					}
					// Send bid-request to frontend
					reply, err := c.ClientBid(ctx, &auction.ClientBidRequest{Id: int32(*port), Amount: int32(bidAmount)})
					if err != nil {
						log.Println("Error: " + err.Error())
					} else {
						log.Printf("Bid was processed: %s", reply.Outcome.String())
					}
				} else if lineSplit[0] == "result" {
					reply, err := c.ClientResult(ctx, &auction.ClientResultRequest{Id: int32(*port)})
					if err != nil {
						log.Println("Error: " + err.Error())
					} else {
						log.Printf("The current highest bid is %d, by bidder %d, IsFinished: %v", reply.Amount, reply.Bidder, reply.IsFinished)
					}
				} else if lineSplit[0] == "auction" {
					reply, err := c.ClientNewAuction(ctx, &auction.ClientNewAuctionRequest{Id: int32(*port)})
					if err != nil {
						log.Println("Error: " + err.Error())
					} else {
						log.Printf("Request was processed. %s", reply.Outcome.String())
					}
					cancel()
				} else {
					log.Println("Did not understand command")
				}
			}
		}
	}
}
