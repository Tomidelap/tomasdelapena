package repositories

import (
	"errors"
	"fmt"
	"sync"

	"clientes/models"
)

var ErrClienteNoEncontrado = errors.New("cliente no encontrado")

type ClientesRepo interface {
	Guardar(cliente models.Cliente) models.Cliente
	BuscarPorID(id string) (models.Cliente, error)
}

type ClientesMemoria struct {
	mu       sync.Mutex
	clientes map[string]models.Cliente
	contador int
}

func NuevoClientesMemoria() *ClientesMemoria {
	return &ClientesMemoria{clientes: make(map[string]models.Cliente)}
}

func (r *ClientesMemoria) Guardar(cliente models.Cliente) models.Cliente {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.contador++
	cliente.ID = fmt.Sprintf("C-%d", r.contador)
	r.clientes[cliente.ID] = cliente
	return cliente
}

func (r *ClientesMemoria) BuscarPorID(id string) (models.Cliente, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cliente, ok := r.clientes[id]
	if !ok {
		return models.Cliente{}, ErrClienteNoEncontrado
	}
	return cliente, nil
}
