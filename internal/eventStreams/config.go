package eventStreams

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var writer *kafka.Producer

func Configure() (w *kafka.Producer, err error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		fmt.Printf("Error creating Kafka producer: %v\n", err)
		return nil, err
	}
	writer = producer
	return writer, nil
}

// func publishOrderEvent(orderData entities.Order) {
// 	// producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
// 	// if err != nil {
// 	// 	fmt.Printf("Error creating Kafka producer: %v\n", err)
// 	// }

// 	// defer producer.Close()

// 	// Delivery report handler for produced messages
// 	go func() {
// 		//p.Events é um canal que emite para aplicação se o evento foi escrito com sucesso ou não
// 		for event := range producer.Events() {
// 			switch eventType := event.(type) {
// 			case *kafka.Message:
// 				if eventType.TopicPartition.Error != nil {
// 					fmt.Printf("Delivery failed: %v\n", eventType.TopicPartition)
// 				} else {
// 					fmt.Printf("Delivered message to %v\n", eventType.TopicPartition)
// 				}
// 			}
// 		}
// 	}()

// 	message, _ := json.Marshal(orderData)

// 	// Produce messages to topic (asynchronously)
// 	topic := "orders-origin-goapi"
// 	producer.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
// 		Value:          []byte(message),
// 	}, nil)

// 	// Wait for message deliveries before shutting down
// 	producer.Flush(15 * 1000)
// }
