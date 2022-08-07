package domain

import (
	"github.com/h2non/bimg"
	"time"
)

type Quality int64

type Image struct {
	Id          int64
	Name        []string
	Quality     []Quality
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type Images struct {
	ObjId *int64
	Items []Image
	Total uint64
	Pages uint64
}

func (q Quality) GetQualities() []Quality {
	return []Quality{100, 75, 50, 25}
}

func (i Image) GetOptions() []bimg.Options {
	var qualities Quality

	options := make([]bimg.Options, 0)

	for _, v := range qualities.GetQualities() {
		options = append(options, bimg.Options{
			Quality: int(v),
		})
	}

	return options
}
