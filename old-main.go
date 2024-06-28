package teste

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/natanebaltezan/orders-service/internal/entities"
	"github.com/natanebaltezan/orders-service/internal/services"
)

var orders = []entities.Order{
	{
		ID:        "123",
		Product:   "Caneta",
		Price:     2.00,
		Priority:  0,
		Status:    "created",
		Timestamp: time.Now(),
	},
	{
		ID:        "124",
		Product:   "Lapis",
		Price:     2.00,
		Priority:  1,
		Status:    "created",
		Timestamp: time.Now(),
	},
}

// c *gin.Context - similar ao req e res
func getOrders(c *gin.Context) {
	// message := consumeOrderEvents()
	// fmt.Println("Message", message)
	// fmt.Printf("%v", reflect.TypeOf(message))
	c.IndentedJSON(http.StatusOK, orders)
}

// TO DO - susbtituir escrita de arquivo para escrita em tópico
func publishOrderEvent(orderData entities.Order) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		fmt.Printf("Error creating Kafka producer: %v\n", err)
	}

	defer producer.Close()

	// Delivery report handler for produced messages
	go func() {
		//p.Events é um canal que emite para aplicação se o evento foi escrito com sucesso ou não
		for event := range producer.Events() {
			switch eventType := event.(type) {
			case *kafka.Message:
				if eventType.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", eventType.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", eventType.TopicPartition)
				}
			}
		}

		fmt.Printf("fechou")
	}()

	message, _ := json.Marshal(orderData)

	// Produce messages to topic (asynchronously)
	topic := "orders-origin-goapi"
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	// Wait for message deliveries before shutting down
	producer.Flush(15 * 1000)
}

// func consumeOrderEvents() string {
// 	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers": "localhost",
// 		"group.id":          "myGroup",
// 		"auto.offset.reset": "earliest",
// 	})

// 	if err != nil {
// 		panic(err)
// 	}
// 	topic := "orders-origin-goapi"
// 	consumer.SubscribeTopics([]string{topic, "^aRegex.*[Tt]opic"}, nil)

// 	// A signal handler or similar could be used to set this to false to break the loop.
// 	run := true

// 	for run {
// 		msg, err := consumer.ReadMessage(time.Second)
// 		fmt.Println("Message", msg)

// 		if err == nil {
// 			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
// 			return string(msg.Value)
// 		} else if !err.(kafka.Error).IsTimeout() {
// 			// The client will automatically try to recover from all errors.
// 			// Timeout is not considered an error because it is raised by
// 			// ReadMessage in absence of messages.
// 			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
// 		}
// 	}

// 	consumer.Close()
// 	// corrigir
// 	return ""
// }

func postOrders(c *gin.Context) {
	var orderData entities.Order

	// BindJSON - preenche no ponteiro de newOrder o objeto/struct com o body recebido na requisição
	if err := c.BindJSON(&orderData); err != nil {
		return
	}

	fmt.Println("POST Orders - Order Data", orderData)
	newOrder := services.BuildOrder(orderData)
	fmt.Println(reflect.TypeOf(newOrder))

	line, _ := json.MarshalIndent(newOrder, "", " ")

	file, err := os.OpenFile("data/orders.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()

	fmt.Println("file", file)
	if err != nil {
		fmt.Println(err)
		return
	}

	n, err := file.Write(line)
	if err != nil {
		fmt.Println(n, err)
	}
	//_ = os.WriteFile("data/orders.json", file, 0644)

	// publishOrderEvent(newOrder)

	c.IndentedJSON(http.StatusCreated, newOrder.ID)
}

func teste() {
	router := gin.Default()
	router.GET("/orders", getOrders)
	router.POST("/order", postOrders)

	router.Run("localhost:8080")
}
