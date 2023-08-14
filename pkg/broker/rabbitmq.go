package broker

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Producer struct {
	conn       *amqp.Connection
	routingKey string
}

func NewProducer(queueURL string) (Producer, error) {
	conn, err := amqp.Dial(queueURL)
	if err != nil {
		return Producer{}, err
	}

	producer := Producer{conn: conn}

	return producer, err
}

func (p *Producer) Close() error {
	return p.conn.Close()
}

func (p *Producer) SendToRabbitMQ(routingKey string, data []byte, header amqp.Table) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Publish(
		"",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Headers:     header,
			Body:        data,
		})
}

func BuildHeader(group, artifactId string) amqp.Table {
	return amqp.Table{
		"schema": fmt.Sprintf("%s/%s", group, artifactId),
	}
}
