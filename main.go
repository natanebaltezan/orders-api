package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type order struct {
	ID        string    `json:"id"`
	Product   string    `json:"product"`
	Price     float64   `json:"price"`
	Priority  int       `json:"priority"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

var orders = []order{
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
	c.IndentedJSON(http.StatusOK, orders)
}

func generateUUID() string {
	id := uuid.NewString()
	return id
}

func buildOrder(orderData order) order {
	orderData.ID = generateUUID()
	orderData.Status = "created"
	orderData.Timestamp = time.Now()

	return orderData
}

func postOrders(c *gin.Context) {
	var orderData order

	// BindJSON - preenche no ponteiro de newOrder o objeto/struct com o body recebido na requisição
	if err := c.BindJSON(&orderData); err != nil {
		return
	}

	fmt.Println("POST Orders - Order Data", orderData)
	newOrder := buildOrder(orderData)
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

	c.IndentedJSON(http.StatusCreated, newOrder.ID)
}

// TO DO - susbtituir escrita de arquivo para escrita em tópico
// func publishOrderEvent(line) {
// 	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
// 	if err != nil {
// 		fmt.Printf("Error creating Kafka producer: %v\n", err)
// 	}

// 	defer producer.Close()

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

// 	// Produce messages to topic (asynchronously)
// 	topic := "orders-origin-goapi"
// 	producer.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
// 		Value:          []byte(line),
// 	}, nil)

// 	// Wait for message deliveries before shutting down
// 	producer.Flush(15 * 1000)
// }

func main() {
	router := gin.Default()
	router.GET("/orders", getOrders)
	router.POST("/order", postOrders)

	router.Run("localhost:8080")
}
