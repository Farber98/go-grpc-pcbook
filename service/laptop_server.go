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

func NewLaptopServer() *LaptopServer {
	return &LaptopServer{}
}

// Unary RPC to create new laptop.
func (s *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Println("Received a create-laptop request with id: %v", laptop.Id)

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
