package application

import (
	"context"

	"github.com/Maxito7/hotel_backend/internal/domain"
	"github.com/Maxito7/hotel_backend/internal/infrastructure/repository"
)

type ContactService struct {
	repo repository.ContactRepository
}

func NewContactService(r repository.ContactRepository) *ContactService {
	return &ContactService{
		repo: r,
	}
}

func (s *ContactService) Create(ctx context.Context, req domain.CreateContactRequest) (int64, error) {
	return s.repo.Create(ctx, req)
}

func (s *ContactService) List(ctx context.Context) ([]domain.Contact, error) {
	return s.repo.List(ctx)
}

func (s *ContactService) UpdateEstado(ctx context.Context, id int64, estado domain.EstadoFormulario) error {
	return s.repo.UpdateEstado(ctx, id, estado)
}
