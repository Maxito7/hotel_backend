package application

import (
	"time"

	"github.com/Maxito7/hotel_backend/internal/domain"
)

type HabitacionService struct {
	repo domain.HabitacionRepository
}

func NewHabitacionService(repo domain.HabitacionRepository) *HabitacionService {
	return &HabitacionService{
		repo: repo,
	}
}

func (s *HabitacionService) GetAllRooms() ([]domain.Habitacion, error) {
	return s.repo.GetAllRooms()
}

func (s *HabitacionService) GetAvailableRooms(fechaEntrada, fechaSalida time.Time) ([]domain.Habitacion, error) {
	return s.repo.GetAvailableRooms(fechaEntrada, fechaSalida)
}

func (s *HabitacionService) GetFechasBloqueadas(desde, hasta time.Time) (*domain.FechasBloqueadas, error) {
	return s.repo.GetFechasBloqueadas(desde, hasta)
}
