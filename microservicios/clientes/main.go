package main

import (
	"log"

	"clientes/controllers"
	"clientes/repositories"
	"clientes/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Cableado: repositorio -> servicio -> controlador
	repo := repositories.NewClientesMemoria()
	service := &services.ClientesService{Repo: repo}
	controller := &controllers.ClientesController{Service: service}

	router := gin.Default()
	router.POST("/clientes", controller.Crear)
	router.GET("/clientes/:id", controller.ObtenerPorID)

	log.Println("Microservicio clientes escuchando en http://localhost:8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
