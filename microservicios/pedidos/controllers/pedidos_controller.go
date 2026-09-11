package controllers

import (
	"net/http"

	"pedidos/services"

	"github.com/gin-gonic/gin"
)

type PedidosController struct {
	Service *services.PedidoService
}

type confirmarPedidoRequest struct {
	ClienteID  string `json:"cliente_id"`
	ProductoID string `json:"producto_id"`
}

func (c *PedidosController) ListarProductos(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"productos": c.Service.ListarProductos()})
}

func (c *PedidosController) ConfirmarPedido(ctx *gin.Context) {
	var req confirmarPedidoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}

	pedido, err := c.Service.ConfirmarPedido(req.ClienteID, req.ProductoID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"mensaje":   "pedido confirmado",
		"pedido_id": pedido.ID,
	})
}
