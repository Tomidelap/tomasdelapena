package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"clientes/repositories"
	"clientes/services"
)

type ClientesController struct {
	service *services.ClientesService
}

func NuevoClientesController(service *services.ClientesService) *ClientesController {
	return &ClientesController{service: service}
}

type crearClienteRequest struct {
	Nombre string `json:"nombre" binding:"required"`
}

func (ctrl *ClientesController) Crear(c *gin.Context) {
	var req crearClienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cliente := ctrl.service.Crear(req.Nombre)
	c.JSON(http.StatusCreated, cliente)
}

func (ctrl *ClientesController) Obtener(c *gin.Context) {
	id := c.Param("id")

	cliente, err := ctrl.service.Buscar(id)
	if err != nil {
		if err == repositories.ErrClienteNoEncontrado {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cliente)
}
