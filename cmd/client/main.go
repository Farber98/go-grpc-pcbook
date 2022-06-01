package main

import (
	"bufio"
	"context"
	"flag"
	"go-grpc-pcbook/pb"
	"go-grpc-pcbook/sample"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient pb.LaptopServiceClient, laptop *pb.Laptop) {
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

func uploadImage(laptopClient pb.LaptopServiceClient, laptopId, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("couldn't open file: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.UploadImage(ctx)
	if err != nil {
		log.Fatal("couldn't upload image: ", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptopId,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("couldn't send image info: ", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("couldn't read chunk from buffer: ", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("couldn't send chunk to server: ", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("couldn't receive response: ", err, stream.RecvMsg(nil)) // get error from sv
	}

	log.Printf("image uploaded with id %s and size %d", res.GetId(), res.GetSize())
}

func testUploadImage(laptopClient pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
	createLaptop(laptopClient, laptop)
	uploadImage(laptopClient, laptop.GetId(), "tmp/laptop.jpg")
}

func testCreateLaptop(laptopClient pb.LaptopServiceClient) {
	createLaptop(laptopClient, sample.NewLaptop())
}

func testSearchLaptop(laptopClient pb.LaptopServiceClient) {
	for i := 0; i < 10; i++ {
		createLaptop(laptopClient, sample.NewLaptop())
	}

	filter := &pb.Filter{
		MaxPrice: 3000,
		MinCores: 4,
		MinGhz:   2.5,
		MinRam:   &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	searchLaptop(laptopClient, filter)
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

	testUploadImage(laptopClient)

}
