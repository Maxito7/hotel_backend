package repository

import (
	"database/sql"
	"fmt"

	"github.com/Maxito7/hotel_backend/internal/domain"
)

type reservaRepository struct {
	db *sql.DB
}

// NewReservaRepository crea una nueva instancia del repositorio de reservas
func NewReservaRepository(db *sql.DB) domain.ReservaRepository {
	return &reservaRepository{db: db}
}

// GetReservaByID obtiene una reserva por su ID con sus habitaciones
func (r *reservaRepository) GetReservaByID(id int) (*domain.Reserva, error) {
	query := `
		SELECT 
			r.reservaid,
			r.cantidadadultos,
			r.cantidadniños,
			r.estado,
			r.clienteid,
			r.subtotal,
			r.descuento,
			r.fechaconfirmacion
		FROM reserva r
		WHERE r.reservaid = $1
	`

	reserva := &domain.Reserva{}
	err := r.db.QueryRow(query, id).Scan(
		&reserva.ID,
		&reserva.CantidadAdultos,
		&reserva.CantidadNinhos,
		&reserva.Estado,
		&reserva.ClienteID,
		&reserva.Subtotal,
		&reserva.Descuento,
		&reserva.FechaConfirmacion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reserva con ID %d no encontrada", id)
		}
		return nil, fmt.Errorf("error al obtener reserva: %w", err)
	}

	// Obtener las habitaciones de la reserva
	habitacionesQuery := `
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
		WHERE rh.reservaid = $1 AND rh.estado = 1
	`

	rows, err := r.db.Query(habitacionesQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener habitaciones de la reserva: %w", err)
	}
	defer rows.Close()

	var habitaciones []domain.ReservaHabitacion
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
			return nil, fmt.Errorf("error al escanear habitación: %w", err)
		}

		habitacion.ID = rh.HabitacionID
		rh.Habitacion = &habitacion
		habitaciones = append(habitaciones, rh)
	}

	reserva.Habitaciones = habitaciones
	return reserva, nil
}

// CreateReserva crea una nueva reserva
func (r *reservaRepository) CreateReserva(reserva *domain.Reserva) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	// Insertar la reserva principal
	query := `
		INSERT INTO reserva (
			cantidadadultos,
			cantidadniños,
			estado,
			clienteid,
			subtotal,
			descuento,
			fechaconfirmacion
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING reservaid
	`

	err = tx.QueryRow(
		query,
		reserva.CantidadAdultos,
		reserva.CantidadNinhos,
		reserva.Estado,
		reserva.ClienteID,
		reserva.Subtotal,
		reserva.Descuento,
		reserva.FechaConfirmacion,
	).Scan(&reserva.ID)

	if err != nil {
		return fmt.Errorf("error al crear reserva: %w", err)
	}

	// Insertar las habitaciones de la reserva
	for i := range reserva.Habitaciones {
		habitacionQuery := `
			INSERT INTO reservaxhabitacion (
				reservaid,
				habitacionid,
				precio,
				fechaentrada,
				fechasalida,
				estado
			) VALUES ($1, $2, $3, $4, $5, $6)
		`

		_, err = tx.Exec(
			habitacionQuery,
			reserva.ID,
			reserva.Habitaciones[i].HabitacionID,
			reserva.Habitaciones[i].Precio,
			reserva.Habitaciones[i].FechaEntrada,
			reserva.Habitaciones[i].FechaSalida,
			1, // Estado activo
		)

		if err != nil {
			return fmt.Errorf("error al crear reserva de habitación: %w", err)
		}

		reserva.Habitaciones[i].ReservaID = reserva.ID
		reserva.Habitaciones[i].Estado = 1
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}

	return nil
}

// UpdateReservaEstado actualiza el estado de una reserva
func (r *reservaRepository) UpdateReservaEstado(id int, estado domain.EstadoReserva) error {
	query := `
		UPDATE reserva 
		SET estado = $1 
		WHERE reservaid = $2
	`

	result, err := r.db.Exec(query, estado, id)
	if err != nil {
		return fmt.Errorf("error al actualizar estado de reserva: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reserva con ID %d no encontrada", id)
	}

	return nil
}

// GetReservasCliente obtiene todas las reservas de un cliente
func (r *reservaRepository) GetReservasCliente(clienteID string) ([]domain.Reserva, error) {
	query := `
		SELECT 
			r.reservaid,
			r.cantidadadultos,
			r.cantidadniños,
			r.estado,
			r.clienteid,
			r.subtotal,
			r.descuento,
			r.fechaconfirmacion
		FROM reserva r
		WHERE r.clienteid = $1
		ORDER BY r.fechaconfirmacion DESC
	`

	rows, err := r.db.Query(query, clienteID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener reservas del cliente: %w", err)
	}
	defer rows.Close()

	var reservas []domain.Reserva
	for rows.Next() {
		var reserva domain.Reserva
		err := rows.Scan(
			&reserva.ID,
			&reserva.CantidadAdultos,
			&reserva.CantidadNinhos,
			&reserva.Estado,
			&reserva.ClienteID,
			&reserva.Subtotal,
			&reserva.Descuento,
			&reserva.FechaConfirmacion,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear reserva: %w", err)
		}

		// Obtener las habitaciones de cada reserva
		habitacionesQuery := `
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
			WHERE rh.reservaid = $1 AND rh.estado = 1
		`

		habRows, err := r.db.Query(habitacionesQuery, reserva.ID)
		if err != nil {
			return nil, fmt.Errorf("error al obtener habitaciones: %w", err)
		}

		var habitaciones []domain.ReservaHabitacion
		for habRows.Next() {
			var rh domain.ReservaHabitacion
			var habitacion domain.Habitacion

			err := habRows.Scan(
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
				habRows.Close()
				return nil, fmt.Errorf("error al escanear habitación: %w", err)
			}

			habitacion.ID = rh.HabitacionID
			rh.Habitacion = &habitacion
			habitaciones = append(habitaciones, rh)
		}
		habRows.Close()

		reserva.Habitaciones = habitaciones
		reservas = append(reservas, reserva)
	}

	return reservas, nil
}
