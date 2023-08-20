package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/max-weis/apicurio-schema-registry-rabbitmq/pkg/broker"
	"github.com/max-weis/apicurio-schema-registry-rabbitmq/pkg/schema"
	"github.com/max-weis/apicurio-schema-registry-rabbitmq/pkg/validator"
	"github.com/streadway/amqp"
	"log"
)

var (
	rabbit     string
	queue      string
	registry   string
	globalId   string
	payload    string
	validation bool
)

func main() {
	ctx := context.Background()

	producer, err := broker.NewProducer(rabbit)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	var user map[string]any
	if err := json.Unmarshal([]byte(payload), &user); err != nil {
		log.Fatalf("Error JSON could not be parsed: %s", err)
	}

	if validation {
		client := schema.NewClient(registry)
		s, err := client.GetSchemaByGlobalId(globalId)
		if err != nil {
			log.Fatalf("Error loading schema: %v", err)
		}

		v := validator.NewAsyncAPIValidator(s, "$.components.schemas.user")
		if err := v.Validate(ctx, user); err != nil {
			log.Fatalf("Error validating schema: %v", err)
		}
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Error JSON could not be created: %v", err)
	}

	header := amqp.Table{"schema": globalId}
	if err := producer.SendToRabbitMQ(queue, bytes, header); err != nil {
		log.Fatalf("Error event could not be send: %s", err)
	}

	log.Printf("User was send successfully: %v", user)
}

func init() {
	flag.StringVar(&rabbit, "rabbit", "amqp://localhost:5672/", "Der Host vom rabbitMQ Server")
	flag.StringVar(&queue, "queue", "user", "Die Queue auf dem die Ereignisse veröffentlicht werden sollen")
	flag.StringVar(&registry, "registry", "http://localhost:8080", "Der Host vom apicurio Server")
	flag.StringVar(&globalId, "globalId", "1", "Die ID der Schema")
	flag.StringVar(&payload, "payload", `{"name":"Max","age":24}`, "Das Ereignis, welches gesendet werden soll")
	flag.BoolVar(&validation, "validation", true, "Schaltet die Validierung aus und ein")
	flag.Parse()
}
