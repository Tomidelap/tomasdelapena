package repositories

import (
	"log"
	"sync"
	"time"

	"pedidos/models"
)

const ttlProductosCache = time.Minute

type ProductosCache struct {
	siguiente   ProductosRepo
	mu          sync.Mutex
	datos       []models.Producto
	vencimiento time.Time
}

func NuevoProductosCache(siguiente ProductosRepo) *ProductosCache {
	return &ProductosCache{siguiente: siguiente}
}

func (c *ProductosCache) Listar() []models.Producto {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.datos != nil && time.Now().Before(c.vencimiento) {
		log.Println("CACHE HIT: productos")
		return c.datos
	}

	log.Println("CACHE MISS: productos")
	c.datos = c.siguiente.Listar()
	c.vencimiento = time.Now().Add(ttlProductosCache)
	return c.datos
}
