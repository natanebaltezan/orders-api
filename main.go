package main

import (
	"fmt"
	"net/http"
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

// c *gin.Context - tipo o req e res
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

	// BindJSON - vai preencher no ponteiro de newOrder o objeto/struct com o body recebido na requisição
	if err := c.BindJSON(&orderData); err != nil {
		return
	}

	fmt.Println("POST Orders - Order Data", orderData)
	newOrder := buildOrder(orderData)
	orders = append(orders, newOrder)
	c.IndentedJSON(http.StatusCreated, newOrder)
}

func main() {
	router := gin.Default()
	router.GET("/orders", getOrders)
	router.POST("/order", postOrders)

	router.Run("localhost:8080")
}
