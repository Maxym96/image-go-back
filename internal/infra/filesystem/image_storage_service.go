package filesystem

import (
	"github.com/h2non/bimg"
	"image-go-back/internal/domain"
	"log"
	"os"
	"path"
)

type ImageStorageService interface {
	SaveImage(image domain.Image, content []byte) error
}

type imageStorageService struct {
	loc string
}

func NewImageStorageService(location string) ImageStorageService {
	return imageStorageService{
		loc: location,
	}
}

func (s imageStorageService) SaveImage(image domain.Image, content []byte) error {
	var locations []string

	for i := range image.Name {
		location := path.Join(s.loc, image.Name[i])
		locations = append(locations, location)
	}
	err := writeFileToStorage(image, locations, content)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func writeFileToStorage(image domain.Image, locations []string, buffer []byte) error {
	dirLocation := path.Dir(locations[0])
	err := os.MkdirAll(dirLocation, os.ModePerm)
	if err != nil {
		log.Print(err)
		return err
	}
	for i, options := range image.GetOptions() {
		var newImage []byte
		newImage, err = bimg.NewImage(buffer).Process(options)
		if err != nil {
			log.Print(err)
			return err
		}

		err = os.WriteFile(locations[i], newImage, os.ModePerm)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}
