package main

import (
	"flag"
	"fmt"
	"go-grpc-pcbook/pb"
	"go-grpc-pcbook/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	serverPort := flag.String("port", "", "server port")
	flag.Parse()
	serverAddress := fmt.Sprintf("0.0.0.0:%s", *serverPort)
	log.Print("starting server at ", serverAddress)
	laptopServer := service.NewLaptopServer(service.NewMemoryLaptopStore(), service.NewDiskImageStore("img"))
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error wiring server: %v", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Error wiring server: %v", err)
	}
}
