package services

import (
	"errors"

	"clientes/models"
	"clientes/repositories"
)

type ClientesService struct {
	Repo repositories.ClientesRepo
}

func (s *ClientesService) CrearCliente(nombre string) (models.Cliente, error) {
	if nombre == "" {
		return models.Cliente{}, errors.New("el nombre es obligatorio")
	}
	return s.Repo.Crear(nombre), nil
}

func (s *ClientesService) ObtenerCliente(id string) (models.Cliente, error) {
	return s.Repo.ObtenerPorID(id)
}
