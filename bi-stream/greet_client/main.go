package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"bi-stream/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	doBiDiStreaming(c)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do a BiDi streaming rpc...")

	// we create a steram by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John 1",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John 2",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John 3",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John 4",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John 5",
			},
		},
	}

	resCh := make(chan struct{})
	// send a bunch of messages to the client
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("sending messages: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive a bunch of messages to the client
	go func() {
		// function to receive a bunch of messages
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("error while receiving: %v", err)
				break
			}

			fmt.Printf("received: %v\n", resp.GetResult())
		}
		close(resCh)
	}()
	<-resCh
}
