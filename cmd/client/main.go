package main

import (
	"context"
	"flag"
	"go-grpc-pcbook/pb"
	"go-grpc-pcbook/sample"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	serverAddress := flag.String("addr", "", "server address")
	flag.Parse()
	log.Print("dial server ", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Couldn't dial server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	laptop := sample.NewLaptop()
	laptop.Id = "3bb88927-ec6b-44e2-aecb-49cbf4eb9c3f"
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	// Set req timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Print("Laptop already exists")
		} else {
			log.Fatal("Couldn't create laptop: ", err)
		}
		return
	}

	log.Printf("Created laptop with id: %s", res.Id)
}
