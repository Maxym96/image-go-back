package database

import (
	"image-go-back/internal/domain"
	"log"
	"time"

	"github.com/upper/db/v4"
)

const ImagesTableName = "images"

type quality int64

type image struct {
	Id          int64      `db:"id,omitempty"`
	ObjId       *int64     `db:"obj_id,omitempty"`
	Name        string     `db:"name"`
	Quality     quality    `db:"quality"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type ImageRepository interface {
	Save(img domain.Image) (domain.Images, error)
	Find(odjId int64, resQuery int64) (domain.Image, error)
	FindById(odjId int64, q int64) (domain.Image, error)
}

type imageRepository struct {
	sess db.Session
	coll db.Collection
}

func NewImageRepository(dbSession db.Session) ImageRepository {
	return imageRepository{
		sess: dbSession,
		coll: dbSession.Collection(ImagesTableName),
	}
}

func (r imageRepository) Save(i domain.Image) (domain.Images, error) {
	var imgsID []int64

	imagesWithId := make([]image, 0)

	images := FromDomainModel(i)
	for _, v := range images {
		err := r.coll.InsertReturning(&v)
		if err != nil {
			log.Printf("imageRepository.Save(): %s", err.Error())
			return domain.Images{}, err
		}
		imagesWithId = append(imagesWithId, v)
	}
	for _, v := range imagesWithId {
		imgsID = append(imgsID, v.Id)
	}

	res := r.sess.SQL().Update(ImagesTableName).Set("obj_id = ?", imgsID[0]).Where(db.Cond{"id IN": imgsID})
	_, err := res.Exec()
	if err != nil {
		log.Printf("imageRepository.Save(): %s", err.Error())
		return domain.Images{}, err
	}
	return mapImagesToDomainCollection(imagesWithId), nil

}

func (r imageRepository) Find(odjId int64, resQuery int64) (domain.Image, error) {
	var i image

	err := r.coll.Find(db.Cond{"obj_id": odjId, "quality": resQuery, "deleted_date": nil}).One(&i)
	if err != nil {
		log.Printf("imageRepository.Find(): %s", err.Error())
		return domain.Image{}, err
	}

	return i.ToDomainModel(), nil
}
func (r imageRepository) FindById(odjId int64, q int64) (domain.Image, error) {
	var i image

	err := r.coll.Find(db.Cond{"obj_id": odjId, "quality": q, "deleted_date": nil}).One(&i)
	if err != nil {
		log.Printf("imageRepository.FindById(): %s", err.Error())
		return domain.Image{}, err
	}

	return i.ToDomainModel(), nil
}

func FromDomainModel(img domain.Image) []image {
	var images []image

	dd := img.DeletedDate
	if dd == nil {
		dd = &time.Time{}
	}
	for i, v := range img.Name {
		createdDate := time.Now()
		images = append(images, image{
			Id:          img.Id,
			Name:        v,
			Quality:     quality(img.Quality[i]),
			CreatedDate: createdDate,
			UpdatedDate: createdDate,
			DeletedDate: dd,
		})
	}
	return images
}

func (i image) ToDomainModel() domain.Image {
	return domain.Image{
		Id:          i.Id,
		Name:        []string{i.Name},
		Quality:     []domain.Quality{domain.Quality(i.Quality)},
		CreatedDate: i.CreatedDate,
		UpdatedDate: i.UpdatedDate,
		DeletedDate: i.DeletedDate,
	}
}

func mapImagesToDomainCollection(images []image) domain.Images {
	var result []domain.Image

	if len(images) == 0 {
		result = make([]domain.Image, 0)
	}

	for _, i := range images {
		result = append(result, i.ToDomainModel())
	}
	res := domain.Images{
		Items: result,
	}

	return res
}
