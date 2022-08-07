package controllers

import (
	"github.com/google/uuid"
	"image-go-back/internal/app"
	"image-go-back/internal/domain"
	appErrors "image-go-back/internal/errors"
	"image-go-back/internal/infra/database"
	"image-go-back/internal/infra/http/resources"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ImageController struct {
	imageService    app.ImageService
	rabbitMQService app.RabbitMQService
}

func NewImageController(s app.ImageService, r app.RabbitMQService) ImageController {
	return ImageController{
		imageService:    s,
		rabbitMQService: r,
	}
}

func (c ImageController) AddImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var quality domain.Quality
		var name []string

		getMsgs, err := c.rabbitMQService.GetBodyFromQueue(database.ImageQueue, true)
		if err != nil {
			err = appErrors.ErrGetBodyFromQueue
			BadRequest(w, err)
			return
		}
		filetype := http.DetectContentType(getMsgs.Body)
		if filetype != "image/jpeg" && filetype != "image/png" {
			err = appErrors.ErrFormatFile
			Forbidden(w, err)
			return
		}
		for _, v := range quality.GetQualities() {
			imageName := uuid.NewString() + "_" + strconv.FormatInt(int64(v), 10) + "." + strings.TrimLeft(filetype, "image/")
			name = append(name, imageName)
		}

		img := domain.Image{
			Name:    name,
			Quality: quality.GetQualities(),
		}

		imgs, err := (c.imageService).Save(img, getMsgs.Body)
		if err != nil {
			log.Printf("ImageController: %s", err)
			BadRequest(w, err)
			return
		}

		imagesDto := resources.NewImageResource(imgs)

		Success(w, imagesDto.Serialize())
	}
}

func (c ImageController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		imgId := r.Context().Value(PathImgKey).(domain.Image)

		var imageDto resources.ImageDto
		Success(w, imageDto.DomainToDto(imgId))
	}
}
