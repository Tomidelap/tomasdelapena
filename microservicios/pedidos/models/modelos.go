package models

type Producto struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Stock  int    `json:"stock"`
}

type Pedido struct {
	ID         string `json:"id"`
	ClienteID  string `json:"cliente_id"`
	ProductoID string `json:"producto_id"`
}

type PedidoConfirmado struct {
	Tipo       string `json:"tipo"`
	PedidoID   string `json:"pedido_id"`
	ClienteID  string `json:"cliente_id"`
	ProductoID string `json:"producto_id"`
}
