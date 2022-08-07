package app

import (
	"github.com/streadway/amqp"
	"image-go-back/internal/infra/database"
)

type RabbitMQService interface {
	SendToQueue(queueName string, body []byte) error
	ReceiveFromQueue(queueName string) error
	GetBodyFromQueue(queueName string, autoAck bool) (amqp.Delivery, error)
}

type rabbitMQService struct {
	repo database.RabbitMQRepository
}

func NewRabbitMQService(r database.RabbitMQRepository) RabbitMQService {
	return &rabbitMQService{
		repo: r,
	}
}

func (r rabbitMQService) SendToQueue(queueName string, body []byte) error {
	return (r.repo).SendToQueue(queueName, body)
}

func (r rabbitMQService) ReceiveFromQueue(queueName string) error {
	return (r.repo).ReceiveFromQueue(queueName)
}
func (r rabbitMQService) GetBodyFromQueue(queueName string, autoAck bool) (amqp.Delivery, error) {
	return (r.repo).GetBodyFromQueue(queueName, autoAck)
}
