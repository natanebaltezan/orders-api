package eventStreams

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/natanebaltezan/orders-service/internal/entities"
)

func PublishOrderEvent(orderData entities.Order) {
	// Delivery report handler for produced messages
	go func() {
		// Problema: tá rodando essa função para cada request
		//p.Events é um canal que emite para aplicação se o evento foi escrito com sucesso ou não
		for event := range writer.Events() {
			switch eventType := event.(type) {
			case *kafka.Message:
				if eventType.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", eventType.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", eventType.TopicPartition)
				}
			}
		}
	}()

	message, _ := json.Marshal(orderData)

	// Produce messages to topic (asynchronously)
	topic := "orders-origin-goapi"
	writer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	// Wait for message deliveries before shutting down
	writer.Flush(15 * 1000)
}
