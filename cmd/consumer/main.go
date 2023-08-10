package main

import (
	"github.com/max-weis/apicurio-schema-registry-rabbitmq/avro"
	"github.com/streadway/amqp"
	"log"
)

const (
	rabbitMQURL = "amqp://localhost:5672/"
	routingKey  = "user.queue"

	group      = "de.adesso.ba-demo.registry"
	artifactId = "user"
)

func main() {
	// Establish a connection to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		routingKey,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for d := range msgs {
		schema, err := avro.GetSchema(group, artifactId)
		if err != nil {
			log.Printf("Failed to get schema by ID: %v", err)
			continue
		}

		message, err := avro.DecodeMessage(d.Body, schema)
		if err != nil {
			log.Printf("Failed to decode message: %v", err)
			continue
		}

		log.Printf("Received message: %v", message)
	}
}
