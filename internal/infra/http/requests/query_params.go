package requests

import (
	appErrors "image-go-back/internal/errors"
	"net/http"
	"strconv"
)

type Findable interface {
	Find(int64) (interface{}, error)
}

func DecodeQueryParams(r *http.Request, queryParam string) (int64, error) {
	var (
		err error
		res uint64
	)

	query := r.URL.Query().Get(queryParam)

	if query != "" {
		res, err = strconv.ParseUint(query, 10, 64)
		if err != nil {
			err = appErrors.ErrParsingParams
			return int64(res), err
		}
	}
	if query == "" {
		err = appErrors.ErrQueryParamsEmpty
		return 0, err
	}

	return int64(res), nil
}
