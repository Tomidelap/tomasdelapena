package repositories

import (
	"errors"
	"fmt"
	"sync"

	"clientes/models"
)

type ClientesRepo interface {
	Crear(nombre string) models.Cliente
	ObtenerPorID(id string) (models.Cliente, error)
}

// ClientesMemoria guarda los clientes en un map (sin base de datos).
type ClientesMemoria struct {
	mu       sync.Mutex
	clientes map[string]models.Cliente
	ultimoID int
}

func NewClientesMemoria() *ClientesMemoria {
	return &ClientesMemoria{clientes: map[string]models.Cliente{}}
}

func (r *ClientesMemoria) Crear(nombre string) models.Cliente {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ultimoID++
	cliente := models.Cliente{ID: fmt.Sprintf("C-%d", r.ultimoID), Nombre: nombre}
	r.clientes[cliente.ID] = cliente
	return cliente
}

func (r *ClientesMemoria) ObtenerPorID(id string) (models.Cliente, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cliente, ok := r.clientes[id]
	if !ok {
		return models.Cliente{}, errors.New("cliente no encontrado")
	}
	return cliente, nil
}
