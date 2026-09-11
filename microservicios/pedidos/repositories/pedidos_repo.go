package repositories

import (
	"fmt"
	"sync"

	"pedidos/models"
)

type PedidosRepo interface {
	Guardar(pedido models.Pedido) models.Pedido
}

type PedidosMemoria struct {
	mu       sync.Mutex
	pedidos  map[string]models.Pedido
	contador int
}

func NuevoPedidosMemoria() *PedidosMemoria {
	return &PedidosMemoria{pedidos: make(map[string]models.Pedido)}
}

func (r *PedidosMemoria) Guardar(pedido models.Pedido) models.Pedido {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.contador++
	pedido.ID = fmt.Sprintf("PED-%d", r.contador)
	r.pedidos[pedido.ID] = pedido
	return pedido
}
