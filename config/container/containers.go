package container

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"image-go-back/config"
	"image-go-back/internal/app"
	"image-go-back/internal/infra/database"
	"image-go-back/internal/infra/filesystem"
	"image-go-back/internal/infra/http/controllers"
	"image-go-back/internal/utils"
	"log"
)

type Container struct {
	Services
	Controllers
}

type Services struct {
	app.ImageService
	app.RabbitMQService
}

type Controllers struct {
	ImageController controllers.ImageController
}

func New(conf config.Configuration) Container {
	sess := getDbSess(conf)
	conn := getRabbitMQConn(conf)

	rabbitMQRepository := database.NewRabbitMQRepository(conn)
	rabbitMQService := app.NewRabbitMQService(rabbitMQRepository)
	imageRepository := database.NewImageRepository(sess)
	imageStorageService := filesystem.NewImageStorageService(conf.FileStorageLocation)
	imageService := app.NewImageService(imageRepository, imageStorageService)
	imageController := controllers.NewImageController(imageService, rabbitMQService)

	return Container{
		Services: Services{
			imageService,
			rabbitMQService,
		},
		Controllers: Controllers{
			imageController,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
func getRabbitMQConn(conf config.Configuration) *amqp.Connection {
	urlString := fmt.Sprintf(
		"amqp://%s:%s@%s/",
		conf.RabbitUser,
		conf.RabbitPassword,
		conf.RabbitHost,
	)
	conn, err := amqp.Dial(urlString)
	if err != nil {
		utils.FailOnError(err, "Unable to create new RabbitMQ connection")
	}
	return conn
}
