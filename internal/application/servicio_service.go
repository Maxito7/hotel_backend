package application

import "github.com/Maxito7/hotel_backend/internal/domain"

type ServicioService struct {
	repo domain.ServicioRepository
}

func NewServicioService(repo domain.ServicioRepository) *ServicioService {
	return &ServicioService{
		repo: repo,
	}
}

func (s *ServicioService) GetAllServices() ([]domain.Servicio, error) {
	return s.repo.GetAllServices()
}
