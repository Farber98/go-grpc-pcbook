package service

import (
	"context"
	"errors"
	"go-grpc-pcbook/pb"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
	Store LaptopStore
}

func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{store}
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
	err := s.Store.Save(laptop)
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

	err := s.Store.Search(stream.Context(), filter, func(laptop *pb.Laptop) error {
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
	return nil
}
