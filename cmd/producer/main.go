package main

import (
	"github.com/max-weis/apicurio-schema-registry-rabbitmq/avro"
	"github.com/max-weis/apicurio-schema-registry-rabbitmq/broker"
	"log"
	"math/rand"
	"time"
)

const (
	group      = "de.adesso.ba-demo.registry"
	artifactId = "user"
)

func main() {
	const rabbitMQURL = "amqp://localhost:5672/"
	const routingKey = "user.queue"

	producer, err := broker.NewProducer(rabbitMQURL)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	for {
		schema, err := avro.GetSchema(group, artifactId)
		if err != nil {
			log.Fatalf("Error parsing Avro schema: %v", err)
		}

		user := map[string]interface{}{"name": "John", "age": rand.Intn(120)}

		data, err := avro.ValidateMessage(user, schema)
		if err != nil {
			log.Fatalf("Error validating message: %v", err)
		}

		if err := producer.SendToRabbitMQ(routingKey, data, broker.BuildHeader(group, artifactId)); err != nil {
			log.Fatalf("Error sending message to RabbitMQ: %v", err)
		}

		log.Printf("User sent successfully: %v", user)
		time.Sleep(3 * time.Second)
	}
}
