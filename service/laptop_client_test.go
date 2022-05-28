package service_test

import (
	"context"
	"go-grpc-pcbook/pb"
	"go-grpc-pcbook/sample"
	"go-grpc-pcbook/serializer"
	"go-grpc-pcbook/service"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	laptopServer, serverAddress := startTestLaptopServer(t)

	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedId := laptop.Id

	// save laptop to store
	req := &pb.CreateLaptopRequest{Laptop: laptop}
	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedId, res.Id)

	// check that laptop was saved to store.
	other, err := laptopServer.Store.Find(laptop.Id)
	require.NoError(t, err)
	require.NotNil(t, other)
	requireSameLaptop(t, laptop, other)
}

func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(service.NewMemoryLaptopStore())

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	listener, err := net.Listen("tcp", ":0") // any random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener) // non blocking call

	return laptopServer, listener.Addr().String()

}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)

}

func requireSameLaptop(t *testing.T, laptop, other *pb.Laptop) {
	// because equal it's not correct, we need to transform and compare it's jsons.

	json1, err := serializer.ProtobufToJson(laptop)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJson(laptop)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}
