package service_test

import (
	"context"
	"go-grpc-pcbook/pb"
	"go-grpc-pcbook/sample"
	"go-grpc-pcbook/service"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerCreateLaptop(t *testing.T) {

	laptopNoId := sample.NewLaptop()
	laptopNoId.Id = ""

	laptopRepeatedId := sample.NewLaptop()
	storeWithRepeatedId := service.NewMemoryLaptopStore()
	err := storeWithRepeatedId.Save(laptopRepeatedId)
	require.Nil(t, err)

	laptopInvalidId := sample.NewLaptop()
	laptopInvalidId.Id = "invalid-id"

	testCases := []struct {
		name   string
		laptop *pb.Laptop
		store  service.LaptopStore
		code   codes.Code
	}{
		{
			name:   "Success with id.",
			laptop: sample.NewLaptop(),
			store:  service.NewMemoryLaptopStore(),
			code:   codes.OK,
		},
		{
			name:   "Success without id.",
			laptop: laptopNoId,
			store:  service.NewMemoryLaptopStore(),
			code:   codes.OK,
		},
		{
			name:   "Repeated UUID.",
			laptop: laptopRepeatedId,
			store:  storeWithRepeatedId,
			code:   codes.AlreadyExists,
		},
		{
			name:   "Invalid UUID.",
			laptop: laptopInvalidId,
			store:  service.NewMemoryLaptopStore(),
			code:   codes.InvalidArgument,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req := &pb.CreateLaptopRequest{
				Laptop: tc.laptop,
			}
			server := service.NewLaptopServer(tc.store)

			res, err := server.CreateLaptop(context.Background(), req)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)
				if len(tc.laptop.Id) > 0 {
					require.Equal(t, tc.laptop.Id, res.Id)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}
		})
	}
}
