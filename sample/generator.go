package sample

import (
	"go-grpc-pcbook/pb"

	"github.com/golang/protobuf/ptypes"
)

func NewKeyboard() *pb.Keyboard {
	return &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
}

func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)
	cores := uint32(randomInt(2, 8))
	threads := uint32(randomInt(int(cores), 12))
	minGhz := randomFloat(2.0, 3.5)
	return &pb.CPU{
		Brand:   brand,
		Name:    name,
		Cores:   cores,
		Threads: threads,
		MinGhz:  minGhz,
		MaxGhz:  randomFloat(minGhz, 5.0),
	}
}

func NewGPU() *pb.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)
	minGhz := randomFloat(2.0, 3.5)
	return &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: randomFloat(minGhz, 5.0),
		Memory: &pb.Memory{
			Value: uint64(randomInt(2, 6)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
}

func NewRAM() *pb.Memory {
	return &pb.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit:  pb.Memory_GIGABYTE,
	}
}

func NewHDD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit:  pb.Memory_TERABYTE,
		},
	}
}

func NewSSD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
}

func NewScreen() *pb.Screen {

	return &pb.Screen{
		SizeInch:   float32(randomFloat(13, 17)),
		Resolution: randomScreenResolution(),
		Panel:      randomScreenPanel(),
		Multitouch: randomBool(),
	}
}

func NewLaptop() *pb.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)
	return &pb.Laptop{
		Id:          randomId(),
		Brand:       brand,
		Name:        name,
		Cpu:         NewCPU(),
		Memory:      NewRAM(),
		Gpu:         []*pb.GPU{NewGPU()},
		Storage:     []*pb.Storage{NewHDD(), NewSSD()},
		Screen:      NewScreen(),
		Keyboard:    NewKeyboard(),
		Weight:      &pb.Laptop_WeightKg{WeightKg: randomFloat(1, 4)},
		Price:       randomFloat(1500, 3000),
		ReleaseYear: int32(randomInt(2018, 2022)),
		UpdatedAt:   ptypes.TimestampNow(),
	}
}

/********* Score *********/
func RandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}
