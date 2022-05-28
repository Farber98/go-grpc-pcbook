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
	// Finds laptop in the store
	Find(id string) (*pb.Laptop, error)
	// Filters laptops from the store. Returns one by one via found func.
	Search(filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

type MemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

type databaseLaptopStore struct {
}

func NewMemoryLaptopStore() *MemoryLaptopStore {
	return &MemoryLaptopStore{data: make(map[string]*pb.Laptop)}
}

func (m *MemoryLaptopStore) Save(laptop *pb.Laptop) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	//deep copy so it can't be modified by its pointer from external world inside our storage.
	other, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	m.data[laptop.Id] = other
	return nil
}

func (m *MemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	laptop, ok := m.data[id]
	if !ok {
		return nil, fmt.Errorf("Given laptop it's not in the store")
	}

	// deep copy
	return deepCopy(laptop)

}

func (m *MemoryLaptopStore) Search(filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, laptop := range m.data {
		if isQualified(filter, laptop) {
			// TODO
		}
	}

	return nil
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("Couldn't copy laptop data: %w", err)
	}
	return other, nil
}

// Determines if a laptop qualifies to be returned by searchLaptop filter.
func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPrice() > filter.GetMaxPrice() || laptop.GetCpu().GetCores() < filter.GetMinCores() || laptop.GetCpu().MinGhz < filter.GetMinGhz() || toBit(laptop.GetMemory()) < toBit(filter.GetMinRam()) {
		return false
	}
	return true
}

// converts memory to bit for comparing purposes.
func toBit(memory *pb.Memory) uint64 {
	val := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return val
	case pb.Memory_BYTE:
		return val << 3
	case pb.Memory_MEGABYTE:
		return val << 23
	case pb.Memory_GIGABYTE:
		return val << 33
	case pb.Memory_TERABYTE:
		return val << 43
	default:
		return 0
	}
}
