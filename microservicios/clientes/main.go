package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"clientes/controllers"
	"clientes/repositories"
	"clientes/services"
)

func main() {
	repo := repositories.NuevoClientesMemoria()
	service := services.NuevoClientesService(repo)
	controller := controllers.NuevoClientesController(service)

	router := gin.Default()
	router.POST("/clientes", controller.Crear)
	router.GET("/clientes/:id", controller.Obtener)

	log.Println("Microservicio de clientes escuchando en http://localhost:8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
