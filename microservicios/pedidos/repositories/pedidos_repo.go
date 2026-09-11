package repositories

import (
	"fmt"
	"sync"

	"pedidos/models"
)

type PedidosRepo interface {
	Guardar(clienteID, productoID string) models.Pedido
}

type PedidosMemoria struct {
	mu       sync.Mutex
	pedidos  map[string]models.Pedido
	ultimoID int
}

func NewPedidosMemoria() *PedidosMemoria {
	return &PedidosMemoria{pedidos: map[string]models.Pedido{}}
}

func (r *PedidosMemoria) Guardar(clienteID, productoID string) models.Pedido {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ultimoID++
	pedido := models.Pedido{
		ID:         fmt.Sprintf("PED-%d", r.ultimoID),
		ClienteID:  clienteID,
		ProductoID: productoID,
	}
	r.pedidos[pedido.ID] = pedido
	return pedido
}
