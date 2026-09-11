package controllers

import (
	"net/http"

	"clientes/services"

	"github.com/gin-gonic/gin"
)

type ClientesController struct {
	Service *services.ClientesService
}

type crearClienteRequest struct {
	Nombre string `json:"nombre"`
}

func (c *ClientesController) Crear(ctx *gin.Context) {
	var req crearClienteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}

	cliente, err := c.Service.CrearCliente(req.Nombre)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, cliente)
}

func (c *ClientesController) ObtenerPorID(ctx *gin.Context) {
	cliente, err := c.Service.ObtenerCliente(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cliente)
}
