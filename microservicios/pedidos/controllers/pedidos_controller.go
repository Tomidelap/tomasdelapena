package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"pedidos/services"
)

type PedidosController struct {
	service *services.PedidoService
}

func NuevoPedidosController(service *services.PedidoService) *PedidosController {
	return &PedidosController{service: service}
}

func (ctrl *PedidosController) ListarProductos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"productos": ctrl.service.ListarProductos()})
}

type confirmarPedidoRequest struct {
	ClienteID  string `json:"cliente_id" binding:"required"`
	ProductoID string `json:"producto_id" binding:"required"`
}

func (ctrl *PedidosController) ConfirmarPedido(c *gin.Context) {
	var req confirmarPedidoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pedido, err := ctrl.service.ConfirmarPedido(req.ClienteID, req.ProductoID)
	if err != nil {
		if errors.Is(err, services.ErrProductoNoEncontrado) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pedido)
}
