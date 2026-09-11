package repositories

import "pedidos/models"

type ProductosRepo interface {
	Listar() []models.Producto
}

type ProductosMemoria struct {
	productos []models.Producto
}

func NuevoProductosMemoria() *ProductosMemoria {
	return &ProductosMemoria{
		productos: []models.Producto{
			{ID: "P-1", Nombre: "Auriculares", Stock: 10},
			{ID: "P-2", Nombre: "Teclado", Stock: 8},
		},
	}
}

func (r *ProductosMemoria) Listar() []models.Producto {
	return r.productos
}
