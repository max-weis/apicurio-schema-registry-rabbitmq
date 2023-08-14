package main

import (
	"context"
	"github.com/max-weis/apicurio-schema-registry-rabbitmq/pkg/broker"
	"github.com/max-weis/apicurio-schema-registry-rabbitmq/pkg/registry"
	"github.com/max-weis/apicurio-schema-registry-rabbitmq/pkg/validate"
	"log"
	"math/rand"
	"time"
)

const (
	group      = "default"
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

	client := registry.NewClient("localhost:8080")

	ctx := context.Background()
	for {
		schema, err := client.GetLatestArtifact(ctx, group, artifactId)
		if err != nil {
			log.Fatalf("Error parsing Avro schema: %v", err)
		}

		user := map[string]interface{}{"name": "John", "age": rand.Intn(120)}

		validator := validate.NewValidator(schema)
		ok, err := validator.Validate(user)
		if err != nil {
			log.Fatalf("Error validating message: %v", err)
		}
		if !ok {
			log.Fatalf("Nachricht entspricht nicht dem Format")
		}

		if err := producer.SendToRabbitMQ(routingKey, user, broker.BuildHeader(group, artifactId)); err != nil {
			log.Fatalf("Error sending message to RabbitMQ: %v", err)
		}

		log.Printf("User sent successfully: %v", user)
		time.Sleep(3 * time.Second)
	}
}
