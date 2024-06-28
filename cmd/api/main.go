package main

import (
	"github.com/gin-gonic/gin"
	"github.com/natanebaltezan/orders-service/internal/controllers"
	"github.com/natanebaltezan/orders-service/internal/eventStreams"
)

func main() {
	writer, _ := eventStreams.Configure()
	//tratar o erro
	defer writer.Close()

	router := gin.Default()
	// router.GET("/orders", getOrders)
	router.POST("/order", controllers.PostOrders)

	router.Run("localhost:8080")

}
