package serializer_test

import (
	"go-grpc-pcbook/pb"
	"go-grpc-pcbook/sample"
	"go-grpc-pcbook/serializer"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatalf("didn't get an error but wanted one: %s", got)
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didn't want one: %v", got)
	}
}

func TestFileSerializer(t *testing.T) {
	t.Run("Writing protobuf to binary file", func(t *testing.T) {
		binaryFile := "../tmp/laptop.bin"

		laptop1 := sample.NewLaptop()
		err := serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)
		assertNoError(t, err)
	})

	t.Run("Reading protobuf from binary file.", func(t *testing.T) {
		binaryFile := "../tmp/laptop.bin"

		laptop1 := &pb.Laptop{}

		err := serializer.ReadProtobufFromBinaryFile(binaryFile, laptop1)
		assertNoError(t, err)
	})

	t.Run("Writing and reading the same from binary file.", func(t *testing.T) {

		binaryFile := "../tmp/laptop.bin"

		laptop1 := sample.NewLaptop()
		laptop2 := &pb.Laptop{}

		serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)
		serializer.ReadProtobufFromBinaryFile(binaryFile, laptop2)

		require.True(t, proto.Equal(laptop1, laptop2))

	})

	t.Run("Writing protobuf to JSON file", func(t *testing.T) {
		binaryFile := "../tmp/laptop.json"

		laptop1 := sample.NewLaptop()
		err := serializer.WriteProtobufToJSONFile(laptop1, binaryFile)
		assertNoError(t, err)
	})

}
