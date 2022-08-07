package database

import (
	"bytes"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	appErrors "image-go-back/internal/errors"
	"log"
	"time"
)

const ImageQueue = "save_image"

type RabbitMQRepository interface {
	SendToQueue(queueName string, body []byte) error
	ReceiveFromQueue(queueName string) error
	GetBodyFromQueue(queueName string, autoAck bool) (amqp.Delivery, error)
}

type rabbitMQRepository struct {
	amqpConn *amqp.Connection
}

func NewRabbitMQRepository(amqpConn *amqp.Connection) RabbitMQRepository {
	return rabbitMQRepository{
		amqpConn: amqpConn,
	}
}
func (r rabbitMQRepository) SendToQueue(queueName string, body []byte) error {
	ch, err := r.amqpConn.Channel()
	if err != nil {
		err = appErrors.ErrDeclareQueue
		return err
	}

	d, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = ch.Publish(
		"",
		d.Name,
		false,
		false,
		amqp.Publishing{
			MessageId:    uuid.NewV4().String(),
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	if err != nil {
		return err
	}
	return err
}
func (r rabbitMQRepository) GetBodyFromQueue(queueName string, autoAck bool) (amqp.Delivery, error) {
	ch, err := r.amqpConn.Channel()
	if err != nil {
		err = appErrors.ErrDeclareQueue
		return amqp.Delivery{}, err
	}
	mess, _, err := ch.Get(queueName, autoAck)
	if err != nil {
		return amqp.Delivery{}, err
	}
	return mess, err
}
func (r rabbitMQRepository) ReceiveFromQueue(queueName string) error {
	ch, err := r.amqpConn.Channel()
	if err != nil {
		return err
	}
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}
	var forever chan struct{}

	go func() {
		for d := range msgs {
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			err = d.Ack(false)
			if err != nil {
				return
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return err
}
