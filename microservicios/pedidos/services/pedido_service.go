package services

import (
	"errors"

	"pedidos/models"
	"pedidos/repositories"
)

type Publisher interface {
	Publicar(evento any) error
}

type PedidoService struct {
	ProductosRepo repositories.ProductosRepo
	PedidosRepo   repositories.PedidosRepo
	Publisher     Publisher
}

func (s *PedidoService) ListarProductos() []models.Producto {
	return s.ProductosRepo.Listar()
}

func (s *PedidoService) ConfirmarPedido(clienteID, productoID string) (models.Pedido, error) {
	if clienteID == "" || productoID == "" {
		return models.Pedido{}, errors.New("cliente_id y producto_id son obligatorios")
	}
	if !s.existeProducto(productoID) {
		return models.Pedido{}, errors.New("producto no encontrado")
	}

	// 1. Se confirma el pedido
	pedido := s.PedidosRepo.Guardar(clienteID, productoID)

	// 2. Se publica el evento para logística
	evento := models.PedidoConfirmado{
		Tipo:       "pedido.confirmado",
		PedidoID:   pedido.ID,
		ClienteID:  pedido.ClienteID,
		ProductoID: pedido.ProductoID,
	}
	if err := s.Publisher.Publicar(evento); err != nil {
		return models.Pedido{}, errors.New("no se pudo publicar el evento: " + err.Error())
	}

	return pedido, nil
}

func (s *PedidoService) existeProducto(id string) bool {
	for _, p := range s.ProductosRepo.Listar() {
		if p.ID == id {
			return true
		}
	}
	return false
}
