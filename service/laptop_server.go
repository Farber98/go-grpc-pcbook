package service

import (
	"bytes"
	"context"
	"errors"
	"go-grpc-pcbook/pb"
	"io"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20

type LaptopServer struct {
	LaptopStore LaptopStore
	ImageStore  ImageStore
}

func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore) *LaptopServer {
	return &LaptopServer{laptopStore, imageStore}
}

// Unary RPC to create new laptop.
func (s *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Println("Received a create-laptop request with id: ", laptop.Id)

	if len(laptop.Id) > 0 {
		// Check valid UUID
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Lapotop ID is not a valid UUID: %v", err)
		}
	} else {
		// Generate uuid
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Couldn't generate laptop UUID: %v", err)
		}
		laptop.Id = id.String() // conver UUID to string format.
	}

	// Supposed heavy processing.
	//time.Sleep(4 * time.Second)

	// Check ctx deadline exceeded before saving to storage.
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("deadline exceeded. Aborting create-laptop req with id %s", laptop.Id)
		return nil, status.Error(codes.DeadlineExceeded, "Deadline exceeded.")
	}

	// Check ctx hasn't been cancelled before saving to storage.
	if ctx.Err() == context.Canceled {
		log.Printf("context cancelled. Aborting create-laptop req with id %s", laptop.Id)
		return nil, status.Error(codes.Canceled, "context cancelled.")
	}

	// Save storage.
	err := s.LaptopStore.Save(laptop)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "Couldn't save laptop to store. UUID already exists: %v", err)
		} else {
			return nil, status.Errorf(codes.Internal, "Couldn't save laptop to store: %v", err)
		}
	}

	log.Printf("Saved laptop with id: %s", laptop.Id)

	// Create response with laptop id and return it.
	return &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}, nil
}

func (s *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Println("Received a search-laptop request with filter: ", filter)

	err := s.LaptopStore.Search(stream.Context(), filter, func(laptop *pb.Laptop) error {
		res := &pb.SearchLaptopResponse{Laptop: laptop}

		err := stream.Send(res)
		if err != nil {
			return err
		}

		log.Println("sent laptop with id: ", laptop.GetId())
		return nil
	})

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}

func (s *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		logError(status.Errorf(codes.Unknown, "couldn't receive image info"))
	}

	laptopId := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetLaptopId()
	log.Println("Received an upload-image request for laptop: ", laptopId)

	laptop, err := s.LaptopStore.Find(laptopId)
	if err != nil {
		logError(status.Errorf(codes.Internal, "couldn't find laptop"))
	}

	if laptop == nil {
		return logError(status.Errorf(codes.InvalidArgument, "laptop not found: %s", laptopId))
	}

	//start receiving image chunks.
	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		log.Println("Receiving data...")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}

		if err != nil {
			return logError(status.Errorf(codes.Unknown, "couldn't receive chunk data: %v", err))
		}

		// Keep track of image size
		chunk := req.GetChunkData()
		size := len(chunk)
		imageSize += size

		// check if image size is greater than the allowed.
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize))
		}

		// write chunk to buffer.
		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "couldn't write chunk data: %v", err))

		}

	}

	// flush buffer to store.
	imageId, err := s.ImageStore.Save(laptopId, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "couldn't flush data to store: %v", err))
	}

	//if image is saved successfully, return image response and clos stream.
	res := &pb.UploadImageResponse{Id: laptopId, Size: uint32(imageSize)}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "couldn't send response: %v", err))
	}

	log.Printf("saved image with id: %s and size %d", imageId, imageSize)
	return nil

}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
