package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/natanebaltezan/orders-service/internal/entities"
	"github.com/natanebaltezan/orders-service/internal/eventStreams"
	"github.com/natanebaltezan/orders-service/internal/services"
)

func PostOrders(c *gin.Context) {
	var orderData entities.Order

	// BindJSON - preenche no ponteiro de newOrder o objeto/struct com o body recebido na requisição
	if err := c.BindJSON(&orderData); err != nil {
		return
	}

	fmt.Println("POST Orders - Order Data", orderData)
	newOrder := services.BuildOrder(orderData)

	// se der erro?
	services.WriteFile(newOrder)
	eventStreams.PublishOrderEvent(newOrder)

	// publishOrderEvent(newOrder)

	c.IndentedJSON(http.StatusCreated, newOrder.ID)
}
