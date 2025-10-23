package repository

import (
	"database/sql"
	"fmt"

	"github.com/Maxito7/hotel_backend/internal/domain"
)

type servicioRepository struct {
	db *sql.DB
}

// NewServicioRepository crea una nueva instancia de servicioRepository
func NewServicioRepository(db *sql.DB) domain.ServicioRepository {
	return &servicioRepository{
		db: db,
	}
}

// GetAllServices implementa domain.ServicioRepository
func (r *servicioRepository) GetAllServices() ([]domain.Servicio, error) {
	query := `
		SELECT 
			service_id,
			name,
			description,
			price
		FROM 
			service
		ORDER BY 
			service_id;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying services: %w", err)
	}
	defer rows.Close()

	var servicios []domain.Servicio
	for rows.Next() {
		var s domain.Servicio
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning service: %w", err)
		}
		servicios = append(servicios, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating services: %w", err)
	}

	return servicios, nil
}
