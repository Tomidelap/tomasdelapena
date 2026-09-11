package repositories

import (
	"log"
	"sync"
	"time"

	"pedidos/models"
)

// ProductosCache envuelve a otro ProductosRepo: si el listado está en caché
// y no venció, lo devuelve; si no, lo pide al siguiente repositorio.
type ProductosCache struct {
	NextRepo ProductosRepo
	TTL      time.Duration

	mu        sync.Mutex
	productos []models.Producto
	vence     time.Time
}

func (r *ProductosCache) Listar() []models.Producto {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.productos != nil && time.Now().Before(r.vence) {
		log.Println("CACHE HIT: productos")
		return r.productos
	}

	log.Println("CACHE MISS: productos")
	r.productos = r.NextRepo.Listar()
	r.vence = time.Now().Add(r.TTL)
	return r.productos
}
