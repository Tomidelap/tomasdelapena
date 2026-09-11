package repositories

import "pedidos/models"

type ProductosRepo interface {
	Listar() []models.Producto
}

// ProductosMemoria tiene el catálogo fijo, igual que el monolito.
type ProductosMemoria struct{}

func (r ProductosMemoria) Listar() []models.Producto {
	return []models.Producto{
		{ID: "P-1", Nombre: "Auriculares", Stock: 10},
		{ID: "P-2", Nombre: "Teclado", Stock: 8},
	}
}
