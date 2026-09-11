package main

import (
	"log"
	"time"

	"pedidos/controllers"
	"pedidos/messaging"
	"pedidos/repositories"
	"pedidos/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Infraestructura: RabbitMQ (mismas credenciales que el compose de CLASE_4)
	publisher, err := messaging.NewRabbitMQPublisher("amqp://user:pass@localhost:5672")
	if err != nil {
		log.Fatalf("Error al conectar con RabbitMQ: %v", err)
	}
	defer publisher.Close()

	// 2. Cableado: la caché envuelve al repositorio de productos
	productosRepo := &repositories.ProductosCache{
		NextRepo: repositories.ProductosMemoria{},
		TTL:      1 * time.Minute,
	}
	service := &services.PedidoService{
		ProductosRepo: productosRepo,
		PedidosRepo:   repositories.NewPedidosMemoria(),
		Publisher:     publisher,
	}
	controller := &controllers.PedidosController{Service: service}

	// 3. Servidor
	router := gin.Default()
	router.GET("/productos", controller.ListarProductos)
	router.POST("/pedidos", controller.ConfirmarPedido)

	log.Println("Microservicio pedidos escuchando en http://localhost:8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
