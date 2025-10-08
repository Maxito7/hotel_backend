package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Maxito7/hotel_backend/internal/domain"
)

type habitacionRepository struct {
	db *sql.DB
}

// NewHabitacionRepository creates a new instance of habitacionRepository
func NewHabitacionRepository(db *sql.DB) domain.HabitacionRepository {
	return &habitacionRepository{
		db: db,
	}
}

// GetAllRooms implements domain.HabitacionRepository
func (r *habitacionRepository) GetAllRooms() ([]domain.Habitacion, error) {
	query := `
		SELECT 
			h.habitacionid,
			h.nombre,
			h.numero,
			h.capacidad,
			h.estado,
			h.descripciongeneral,
			t.tipohabitacionid,
			t.titulo,
			t.descripcion,
			t.capacidadadultos,
			t.capacidadninhos,
			t.cantidadcamas,
			t.precio
		FROM 
			habitacion h
		INNER JOIN 
			tipohabitacion t ON h.tipohabitacionid = t.tipohabitacionid
		ORDER BY 
			h.habitacionid;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	var habitaciones []domain.Habitacion
	for rows.Next() {
		var h domain.Habitacion
		err := rows.Scan(
			&h.ID,
			&h.Nombre,
			&h.Numero,
			&h.Capacidad,
			&h.Estado,
			&h.DescripcionGeneral,
			&h.TipoHabitacion.ID,
			&h.TipoHabitacion.Titulo,
			&h.TipoHabitacion.Descripcion,
			&h.TipoHabitacion.CapacidadAdultos,
			&h.TipoHabitacion.CapacidadNinhos,
			&h.TipoHabitacion.CantidadCamas,
			&h.TipoHabitacion.Precio,
		)
		if err != nil {
			return nil, err
		}
		habitaciones = append(habitaciones, h)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return habitaciones, nil
}

// GetAvailableRooms implements domain.HabitacionRepository
func (r *habitacionRepository) GetAvailableRooms(fechaEntrada, fechaSalida time.Time) ([]domain.Habitacion, error) {
	query := `
		SELECT DISTINCT 
			h.habitacionid,
			h.nombre,
			h.numero,
			h.capacidad,
			h.estado,
			h.descripciongeneral,
			t.tipohabitacionid,
			t.titulo,
			t.descripcion,
			t.capacidadadultos,
			t.capacidadninhos,
			t.cantidadcamas,
			t.precio,
			h.mediaid
		FROM 
			habitacion h
		INNER JOIN 
			tipohabitacion t ON h.tipohabitacionid = t.tipohabitacionid
		WHERE 
			h.estado = 'Disponible'
			AND NOT EXISTS (
				SELECT 1 FROM reservaxhabitacion rh
				WHERE rh.habitacionid = h.habitacionid
				AND rh.estado = 1
				AND (
					(rh.fechaentrada <= $1 AND rh.fechasalida >= $1)
					OR (rh.fechaentrada <= $2 AND rh.fechasalida >= $2)
					OR (rh.fechaentrada >= $1 AND rh.fechasalida <= $2)
				)
			)
		ORDER BY 
			h.habitacionid;`

	rows, err := r.db.Query(query, fechaEntrada, fechaSalida)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habitaciones []domain.Habitacion
	for rows.Next() {
		var h domain.Habitacion
		err := rows.Scan(
			&h.ID,
			&h.Nombre,
			&h.Numero,
			&h.Capacidad,
			&h.Estado,
			&h.DescripcionGeneral,
			&h.TipoHabitacion.ID,
			&h.TipoHabitacion.Titulo,
			&h.TipoHabitacion.Descripcion,
			&h.TipoHabitacion.CapacidadAdultos,
			&h.TipoHabitacion.CapacidadNinhos,
			&h.TipoHabitacion.CantidadCamas,
			&h.TipoHabitacion.Precio,
			&h.MediaID,
		)
		if err != nil {
			return nil, err
		}
		habitaciones = append(habitaciones, h)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return habitaciones, nil
}
