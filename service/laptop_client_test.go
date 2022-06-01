package service_test

import (
	"bufio"
	"context"
	"fmt"
	"go-grpc-pcbook/pb"
	"go-grpc-pcbook/sample"
	"go-grpc-pcbook/serializer"
	"go-grpc-pcbook/service"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	laptopStore := service.NewMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore, nil)

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
	other, err := laptopStore.Find(laptop.Id)
	require.NoError(t, err)
	require.NotNil(t, other)
	requireSameLaptop(t, laptop, other)
}

func TestClientUploadImage(t *testing.T) {
	testImageFolder := "../tmp"

	laptopStore := service.NewMemoryLaptopStore()
	imageStore := service.NewDiskImageStore(testImageFolder)

	laptop := sample.NewLaptop()
	err := laptopStore.Save(laptop)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopStore, imageStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	imagePath := fmt.Sprintf("%s/laptop.jpg", testImageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := laptopClient.UploadImage(context.Background())
	require.NoError(t, err)

	imageType := filepath.Ext(imagePath)
	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptop.GetId(),
				ImageType: imageType,
			},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		size += n

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetId())
	require.EqualValues(t, size, res.GetSize())

	savedImagePath := fmt.Sprintf("%s/%s%s", testImageFolder, res.GetId(), imageType)
	require.FileExists(t, savedImagePath)
	require.NoError(t, os.Remove(savedImagePath))

}

func startTestLaptopServer(t *testing.T, laptopStore service.LaptopStore, imageStore service.ImageStore) string {
	laptopServer := service.NewLaptopServer(laptopStore, imageStore)

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	listener, err := net.Listen("tcp", ":0") // any random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener) // non blocking call

	return listener.Addr().String()

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

func TestClientSearchLaptop(t *testing.T) {
	filter := &pb.Filter{
		MaxPrice: 2000,
		MinCores: 4,
		MinGhz:   2.2,
		MinRam:   &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	store := service.NewMemoryLaptopStore()
	expectedIds := make(map[string]bool)

	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()

		switch i {
		case 0:
			laptop.Price = 2500
		case 1:
			laptop.Cpu.Cores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Memory = &pb.Memory{Value: 4096, Unit: pb.Memory_MEGABYTE}
		case 4:
			laptop.Price = 2000
			laptop.Cpu.Cores = 4
			laptop.Cpu.MinGhz = 2.2
			laptop.Memory = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedIds[laptop.Id] = true
		case 5:
			laptop.Price = 1000
			laptop.Cpu.Cores = 8
			laptop.Cpu.MinGhz = 2.9
			laptop.Memory = &pb.Memory{Value: 32, Unit: pb.Memory_GIGABYTE}
			expectedIds[laptop.Id] = true
		}
		err := store.Save(laptop)
		require.NoError(t, err)
	}

	serverAddress := startTestLaptopServer(t, store, nil)

	laptopClient := newTestLaptopClient(t, serverAddress)

	req := &pb.SearchLaptopRequest{Filter: filter}

	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)

	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.Contains(t, expectedIds, res.GetLaptop().GetId())
		found++
	}
	require.Equal(t, len(expectedIds), found)
}
