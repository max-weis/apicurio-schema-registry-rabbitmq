package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/max-weis/apicurio-schema-registry-rabbitmq/pkg/schema"
	"github.com/streadway/amqp"
	"log"
)

var (
	rabbit   string
	queue    string
	registry string
)

func main() {
	// Establish a connection to RabbitMQ
	conn, err := amqp.Dial(rabbit)
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
		queue,
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

	client := schema.NewClient(registry)

	for d := range msgs {
		ctx := context.Background()
		s, err := client.GetSchemaByGlobalId(d.Headers["schema"].(string))
		if err != nil {
			log.Fatalf("Failed to get schema: %v", err)
		}

		var user map[string]any
		if err = json.Unmarshal(d.Body, &user); err != nil {
			log.Fatalf("Failed to unmarshal json: %v", err)
		}

		validator := schema.NewAvroValidator(s)
		ok, err := validator.Validate(ctx, user)
		if err != nil {
			log.Printf("Failed to validate schema: %v: %v", user, err)
		}

		if !ok {
			log.Println("Message could not be validated")
		}

		log.Printf("Received message: %#v", user)
	}
}

func init() {
	flag.StringVar(&rabbit, "rabbit", "amqp://localhost:5672/", "Der Host vom rabbitMQ Server")
	flag.StringVar(&queue, "queue", "user", "Die Queue auf dem die Ereignisse veröffentlicht werden sollen")
	flag.StringVar(&registry, "registry", "http://localhost:8080", "Der Host vom apicurio Server")
	flag.Parse()
}
