package app

import (
	"image-go-back/internal/domain"
	"image-go-back/internal/infra/database"
	"image-go-back/internal/infra/filesystem"
	"log"
)

type ImageService interface {
	Save(image domain.Image, content []byte) (domain.Images, error)
	FindById(objId int64, q int64) (domain.Image, error)
	Find(int64, int64) (interface{}, error)
}

type imageService struct {
	repo    database.ImageRepository
	filesys filesystem.ImageStorageService
}

func NewImageService(r database.ImageRepository, s filesystem.ImageStorageService) ImageService {
	return &imageService{
		repo:    r,
		filesys: s,
	}
}

func (s imageService) Save(image domain.Image, content []byte) (domain.Images, error) {
	err := s.filesys.SaveImage(image, content)
	if err != nil {
		log.Print(err)
		return domain.Images{}, err
	}

	imgs, err := s.repo.Save(image)
	if err != nil {
		log.Print(err)
		return domain.Images{}, err
	}

	return imgs, nil
}

func (s imageService) FindById(objId int64, q int64) (domain.Image, error) {
	return s.repo.FindById(objId, q)
}

func (s imageService) Find(objId int64, resQuery int64) (interface{}, error) {
	return s.repo.Find(objId, resQuery)
}
