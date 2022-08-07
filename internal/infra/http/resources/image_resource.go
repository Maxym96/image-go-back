package resources

import (
	"fmt"
	"image-go-back/internal/domain"
)

type ImageDto struct {
	Quality int64  `json:"quality"`
	Link    string `json:"link"`
}
type ImagesDto struct {
	ObjId *int64     `json:"obj_id"`
	Items []ImageDto `json:"items"`
	Pages uint64     `json:"pages"`
	Total uint64     `json:"total"`
}

type ImagesResource struct {
	images domain.Images
}

func NewImageResource(images domain.Images) ImagesResource {
	return ImagesResource{images: images}
}

func (ir ImagesResource) Serialize() ImagesDto {
	var (
		result []ImageDto
		imgDto ImageDto
	)

	if len(ir.images.Items) == 0 {
		result = make([]ImageDto, 0)
	}

	for _, i := range ir.images.Items {
		result = append(result, imgDto.DomainToDto(i))
	}

	res := ImagesDto{
		ObjId: &ir.images.Items[0].Id,
		Items: result,
		Total: ir.images.Total,
		Pages: ir.images.Pages,
	}

	return res
}

func (d ImageDto) DomainToDto(i domain.Image) ImageDto {
	link := fmt.Sprintf("/static/%s", i.Name[0])
	d.Quality = int64(i.Quality[0])
	d.Link = link
	return d
}
