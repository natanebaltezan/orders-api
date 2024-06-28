package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/natanebaltezan/orders-service/internal/entities"
)

func generateUUID() string {
	id := uuid.NewString()
	return id
}

func BuildOrder(orderData entities.Order) entities.Order {
	orderData.ID = generateUUID()
	orderData.Status = "created"
	orderData.Timestamp = time.Now()

	return orderData
}

func WriteFile(newOrder entities.Order) {
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
}
