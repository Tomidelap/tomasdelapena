package services

import (
	"clientes/models"
	"clientes/repositories"
)

type ClientesService struct {
	repo repositories.ClientesRepo
}

func NuevoClientesService(repo repositories.ClientesRepo) *ClientesService {
	return &ClientesService{repo: repo}
}

func (s *ClientesService) Crear(nombre string) models.Cliente {
	return s.repo.Guardar(models.Cliente{Nombre: nombre})
}

func (s *ClientesService) Buscar(id string) (models.Cliente, error) {
	return s.repo.BuscarPorID(id)
}
