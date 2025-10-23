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
		INSERT INTO reservation_room (
			reservation_id,
			room_id,
			price,
			check_in_date,
			check_out_date,
			status
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
			rh.reservation_id,
			rh.room_id,
			rh.price,
			rh.check_in_date,
			rh.check_out_date,
			rh.status,
			h.name,
			h.capacity,
			h.number
		FROM reservation_room rh
		INNER JOIN room h ON h.room_id = rh.room_id
		WHERE rh.reservation_id = $1
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
		UPDATE reservation_room 
		SET status = $1 
		WHERE reservation_id = $2 AND room_id = $3
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
		FROM reservation_room rh
		INNER JOIN reservation r ON r.reservation_id = rh.reservation_id
		WHERE rh.room_id = $1 
		AND rh.status = 1
		AND r.status NOT IN ('Cancelada')
		AND (
			(rh.check_in_date < $3 AND rh.check_out_date > $2)
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
			rh.status,
			h.name,
			h.capacity,
			h.number
		FROM reservation_room rh
		INNER JOIN room h ON h.room_id = rh.room_id
		INNER JOIN reservation r ON r.reservation_id = rh.reservation_id
		WHERE rh.status = 1
		AND r.status NOT IN ('Cancelada')
		AND (
			(rh.check_in_date BETWEEN $1 AND $2)
			OR (rh.check_out_date BETWEEN $1 AND $2)
			OR (rh.check_in_date < $1 AND rh.check_out_date > $2)
		)
		ORDER BY rh.check_in_date
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
