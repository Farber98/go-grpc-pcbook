package service

import (
	"errors"
	"fmt"
	"go-grpc-pcbook/pb"
	"sync"

	"github.com/jinzhu/copier"
)

var (
	ErrAlreadyExists = errors.New("UUID already exists.")
)

type LaptopStore interface {
	// Saves laptop to the store
	Save(laptop *pb.Laptop) error
}

type MemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

type databaseLaptopStore struct {
}

func (m *MemoryLaptopStore) NewMemoryLaptopStore(laptop *pb.Laptop) *MemoryLaptopStore {
	return &MemoryLaptopStore{data: make(map[string]*pb.Laptop)}
}

func (m *MemoryLaptopStore) Save(laptop *pb.Laptop) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	//deep copy so it can't be modified by its pointer from external world inside our storage.
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("Couldn't copy laptop data: %w", err)
	}

	m.data[laptop.Id] = other
	return nil
}
