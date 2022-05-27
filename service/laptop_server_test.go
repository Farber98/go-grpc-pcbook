package service_test

import (
	"go-grpc-pcbook/pb"
	"go-grpc-pcbook/service"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestXxx(t *testing.T) {

	testCases := []struct {
		name   string
		laptop *pb.Laptop
		store  service.MemoryLaptopStore
		code   codes.Code
	}{}
}
