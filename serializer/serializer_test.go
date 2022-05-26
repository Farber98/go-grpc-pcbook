package serializer_test

import (
	"go-grpc-pcbook/sample"
	"go-grpc-pcbook/serializer"
	"testing"
)

func TestFileSerializer(t *testing.T) {
	t.Run("Writing protobuf to binary file", func(t *testing.T) {
		binaryFile := "../tmp/laptop.bin"

		laptop1 := sample.NewLaptop()
		err := serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)
		assertNoError(t, err)
	})
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didn't want one: %v", got)
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
