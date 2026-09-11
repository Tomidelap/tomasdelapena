package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"pedidos/controllers"
	"pedidos/messaging"
	"pedidos/repositories"
	"pedidos/services"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	publisher, err := messaging.NuevoRabbitMQPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("no se pudo conectar a RabbitMQ: %v", err)
	}

	productosRepo := repositories.NuevoProductosMemoria()
	productosCache := repositories.NuevoProductosCache(productosRepo)
	pedidosRepo := repositories.NuevoPedidosMemoria()

	service := services.NuevoPedidoService(productosCache, pedidosRepo, publisher)
	controller := controllers.NuevoPedidosController(service)

	router := gin.Default()
	router.GET("/productos", controller.ListarProductos)
	router.POST("/pedidos", controller.ConfirmarPedido)

	log.Println("Microservicio de pedidos escuchando en http://localhost:8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
