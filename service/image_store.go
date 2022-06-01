package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type ImageStore interface {
	Save(laptopId string, imageType string, imageData bytes.Buffer) (string, error)
}

type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*ImageInfo
}

type ImageInfo struct {
	LaptopId string
	Type     string
	Path     string
}

func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{imageFolder: imageFolder, images: make(map[string]*ImageInfo)}
}

func (d *DiskImageStore) Save(laptopId string, imageType string, imageData bytes.Buffer) (string, error) {
	imageId, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("Couldn't create image id: %v", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", d.imageFolder, imageId, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("Couldn't create image file: %v", err)
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("Couldn't write image to file: %v", err)
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.images[imageId.String()] = &ImageInfo{LaptopId: laptopId, Type: imageType, Path: imagePath}

	return imageId.String(), nil
}
