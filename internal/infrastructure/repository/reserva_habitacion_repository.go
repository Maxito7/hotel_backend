package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Maxito7/hotel_backend/internal/domain"
)

type reservaHabitacionRepository struct {
	db *sql.DB
}

// NewReservaHabitacionRepository crea una nueva instancia del repositorio
func NewReservaHabitacionRepository(db *sql.DB) domain.ReservaHabitacionRepository {
	return &reservaHabitacionRepository{db: db}
}

// CreateReservaHabitacion crea una nueva reserva de habitación
func (r *reservaHabitacionRepository) CreateReservaHabitacion(reservaHabitacion *domain.ReservaHabitacion) error {
	query := `
		INSERT INTO reservaxhabitacion (
			reservaid,
			habitacionid,
			precio,
			fechaentrada,
			fechasalida,
			estado
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		query,
		reservaHabitacion.ReservaID,
		reservaHabitacion.HabitacionID,
		reservaHabitacion.Precio,
		reservaHabitacion.FechaEntrada,
		reservaHabitacion.FechaSalida,
		reservaHabitacion.Estado,
	)

	if err != nil {
		return fmt.Errorf("error al crear reserva de habitación: %w", err)
	}

	return nil
}

// GetReservaHabitacionesByReservaID obtiene todas las habitaciones de una reserva
func (r *reservaHabitacionRepository) GetReservaHabitacionesByReservaID(reservaID int) ([]domain.ReservaHabitacion, error) {
	query := `
		SELECT 
			rh.reservaid,
			rh.habitacionid,
			rh.precio,
			rh.fechaentrada,
			rh.fechasalida,
			rh.estado,
			h.nombre,
			h.capacidad,
			h.numero
		FROM reservaxhabitacion rh
		INNER JOIN habitacion h ON h.habitacionid = rh.habitacionid
		WHERE rh.reservaid = $1
	`

	rows, err := r.db.Query(query, reservaID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener habitaciones de la reserva: %w", err)
	}
	defer rows.Close()

	var reservasHabitacion []domain.ReservaHabitacion
	for rows.Next() {
		var rh domain.ReservaHabitacion
		var habitacion domain.Habitacion

		err := rows.Scan(
			&rh.ReservaID,
			&rh.HabitacionID,
			&rh.Precio,
			&rh.FechaEntrada,
			&rh.FechaSalida,
			&rh.Estado,
			&habitacion.Nombre,
			&habitacion.Capacidad,
			&habitacion.Numero,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear reserva de habitación: %w", err)
		}

		habitacion.ID = rh.HabitacionID
		rh.Habitacion = &habitacion
		reservasHabitacion = append(reservasHabitacion, rh)
	}

	return reservasHabitacion, nil
}

// UpdateReservaHabitacionEstado actualiza el estado de una reserva de habitación
func (r *reservaHabitacionRepository) UpdateReservaHabitacionEstado(reservaID, habitacionID int, estado int) error {
	query := `
		UPDATE reservaxhabitacion 
		SET estado = $1 
		WHERE reservaid = $2 AND habitacionid = $3
	`

	result, err := r.db.Exec(query, estado, reservaID, habitacionID)
	if err != nil {
		return fmt.Errorf("error al actualizar estado: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reserva de habitación no encontrada")
	}

	return nil
}

// VerificarDisponibilidad verifica si una habitación está disponible para las fechas dadas
func (r *reservaHabitacionRepository) VerificarDisponibilidad(habitacionID int, fechaEntrada, fechaSalida time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM reservaxhabitacion rh
		INNER JOIN reserva r ON r.reservaid = rh.reservaid
		WHERE rh.habitacionid = $1 
		AND rh.estado = 1
		AND r.estado NOT IN ('Cancelada')
		AND (
			(rh.fechaentrada < $3 AND rh.fechasalida > $2)
		)
	`

	var count int
	err := r.db.QueryRow(query, habitacionID, fechaEntrada, fechaSalida).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error al verificar disponibilidad: %w", err)
	}

	// Si count es 0, la habitación está disponible
	return count == 0, nil
}

// GetReservasEnRango obtiene todas las reservas activas en un rango de fechas
func (r *reservaHabitacionRepository) GetReservasEnRango(fechaInicio, fechaFin time.Time) ([]domain.ReservaHabitacion, error) {
	query := `
		SELECT 
			rh.reservaid,
			rh.habitacionid,
			rh.precio,
			rh.fechaentrada,
			rh.fechasalida,
			rh.estado,
			h.nombre,
			h.capacidad,
			h.numero
		FROM reservaxhabitacion rh
		INNER JOIN habitacion h ON h.habitacionid = rh.habitacionid
		INNER JOIN reserva r ON r.reservaid = rh.reservaid
		WHERE rh.estado = 1
		AND r.estado NOT IN ('Cancelada')
		AND (
			(rh.fechaentrada BETWEEN $1 AND $2)
			OR (rh.fechasalida BETWEEN $1 AND $2)
			OR (rh.fechaentrada < $1 AND rh.fechasalida > $2)
		)
		ORDER BY rh.fechaentrada
	`

	rows, err := r.db.Query(query, fechaInicio, fechaFin)
	if err != nil {
		return nil, fmt.Errorf("error al obtener reservas en rango: %w", err)
	}
	defer rows.Close()

	var reservasHabitacion []domain.ReservaHabitacion
	for rows.Next() {
		var rh domain.ReservaHabitacion
		var habitacion domain.Habitacion

		err := rows.Scan(
			&rh.ReservaID,
			&rh.HabitacionID,
			&rh.Precio,
			&rh.FechaEntrada,
			&rh.FechaSalida,
			&rh.Estado,
			&habitacion.Nombre,
			&habitacion.Capacidad,
			&habitacion.Numero,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear reserva: %w", err)
		}

		habitacion.ID = rh.HabitacionID
		rh.Habitacion = &habitacion
		reservasHabitacion = append(reservasHabitacion, rh)
	}

	return reservasHabitacion, nil
}
