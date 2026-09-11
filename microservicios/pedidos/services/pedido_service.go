package services

import (
	"errors"
	"fmt"

	"pedidos/messaging"
	"pedidos/models"
	"pedidos/repositories"
)

var ErrProductoNoEncontrado = errors.New("producto no encontrado")

type PedidoService struct {
	productos repositories.ProductosRepo
	pedidos   repositories.PedidosRepo
	publisher messaging.Publisher
}

func NuevoPedidoService(productos repositories.ProductosRepo, pedidos repositories.PedidosRepo, publisher messaging.Publisher) *PedidoService {
	return &PedidoService{productos: productos, pedidos: pedidos, publisher: publisher}
}

func (s *PedidoService) ListarProductos() []models.Producto {
	return s.productos.Listar()
}

func (s *PedidoService) ConfirmarPedido(clienteID, productoID string) (models.Pedido, error) {
	if !s.existeProducto(productoID) {
		return models.Pedido{}, ErrProductoNoEncontrado
	}

	pedido := s.pedidos.Guardar(models.Pedido{ClienteID: clienteID, ProductoID: productoID})

	evento := models.PedidoConfirmado{
		Tipo:       "pedido.confirmado",
		PedidoID:   pedido.ID,
		ClienteID:  pedido.ClienteID,
		ProductoID: pedido.ProductoID,
	}
	if err := s.publisher.PublicarPedidoConfirmado(evento); err != nil {
		return models.Pedido{}, fmt.Errorf("pedido guardado pero fallo al publicar el evento: %w", err)
	}

	return pedido, nil
}

func (s *PedidoService) existeProducto(id string) bool {
	for _, p := range s.productos.Listar() {
		if p.ID == id {
			return true
		}
	}
	return false
}
