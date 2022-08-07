package middlewares

import (
	"context"
	appErrors "image-go-back/internal/errors"
	"image-go-back/internal/infra/http/controllers"
	"image-go-back/internal/infra/http/requests"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/upper/db/v4"
)

type Findable interface {
	Find(int64, int64) (interface{}, error)
}

func PathObject(pathKey string, ctxKey string, queryParam string, service Findable) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseUint(chi.URLParam(r, pathKey), 10, 64)
			if err != nil {
				err = appErrors.ErrInvalidPathKey
				log.Print(err)
				controllers.BadRequest(w, err)
				return
			}

			q, err := requests.DecodeQueryParams(r, queryParam)
			if err != nil {
				log.Print(err)
				controllers.BadRequest(w, err)
				return
			}

			obj, err := service.Find(int64(id), q)
			if err != nil {
				if err == db.ErrNoMoreRows {
					log.Print(err)
					controllers.NotFound(w, err)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ctxKey, obj)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
