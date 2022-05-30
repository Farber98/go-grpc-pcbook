package main

import (
	"context"
	"flag"
	"go-grpc-pcbook/pb"
	"go-grpc-pcbook/sample"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient pb.LaptopServiceClient) {

	laptop := sample.NewLaptop()
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

func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
	log.Println("search filter: ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{Filter: filter}

	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("couldn't search laptop: ", err)
	}
	count := 1
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("couldn't receive response: ", err)
		}
		log.Printf(">>>>>> Laptop %d <<<<<<", count)
		log.Println("brand: ", res.Laptop.GetBrand())
		log.Println("name: ", res.Laptop.GetName())
		log.Println("cores: ", res.Laptop.GetCpu().GetCores())
		log.Println("min freq: ", res.Laptop.GetCpu().GetMinGhz())
		log.Println("ram: ", res.Laptop.GetMemory().GetValue(), res.Laptop.GetMemory().GetUnit())
		log.Println("price: $", res.Laptop.GetPrice())

		count++
	}
}

func main() {
	serverAddress := flag.String("addr", "", "server address")
	flag.Parse()
	log.Print("dial server ", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Couldn't dial server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	for i := 0; i < 10; i++ {
		createLaptop(laptopClient)
	}

	filter := &pb.Filter{
		MaxPrice: 3000,
		MinCores: 4,
		MinGhz:   2.5,
		MinRam:   &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	searchLaptop(laptopClient, filter)
}
