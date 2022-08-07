package middlewares

import (
	"image-go-back/internal/app"
	appErrors "image-go-back/internal/errors"
	"image-go-back/internal/infra/database"
	"image-go-back/internal/infra/http/controllers"
	"io/ioutil"
	"net/http"
)

func SaveImage(rmq app.RabbitMQService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			buff, err := ioutil.ReadAll(r.Body)
			if err != nil {
				err = appErrors.ErrReadAllBody
				controllers.BadRequest(w, err)
			}
			err = rmq.SendToQueue(database.ImageQueue, buff)
			if err != nil {
				err = appErrors.ErrSendToQueue
				controllers.BadRequest(w, err)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
