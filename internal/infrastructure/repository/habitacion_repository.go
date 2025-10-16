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

// GetDisponibilidadFechas implementa domain.HabitacionRepository
func (r *habitacionRepository) GetDisponibilidadFechas(desde, hasta time.Time) ([]domain.DisponibilidadFecha, error) {
	query := `
		WITH RECURSIVE fechas AS (
			SELECT date(cast($1 as timestamp)) as fecha
			UNION ALL
			SELECT fecha + interval '1 day'
			FROM fechas
			WHERE fecha < date(cast($2 as timestamp))
		),
		habitaciones_totales AS (
			SELECT COUNT(*) as total
			FROM habitacion
			WHERE estado = 'Disponible'
		),
		habitaciones_ocupadas AS (
			SELECT date(f.fecha) as fecha, 
				   COUNT(DISTINCT rh.habitacionid) as ocupadas
			FROM fechas f
			LEFT JOIN reservaxhabitacion rh ON 
				date(f.fecha) BETWEEN date(rh.fechaentrada) AND date(rh.fechasalida)
			LEFT JOIN reserva r ON r.reservaid = rh.reservaid
				AND rh.estado = 1
				AND r.estado = 'Confirmada'
			GROUP BY date(f.fecha)
		)
		SELECT 
			f.fecha,
			CASE 
				WHEN (ht.total - COALESCE(ho.ocupadas, 0)) > 0 THEN true 
				ELSE false 
			END as disponible,
			(ht.total - COALESCE(ho.ocupadas, 0)) as habitaciones_disponibles
		FROM fechas f
		CROSS JOIN habitaciones_totales ht
		LEFT JOIN habitaciones_ocupadas ho ON date(f.fecha) = date(ho.fecha)
		ORDER BY f.fecha;`

	rows, err := r.db.Query(query, desde, hasta)
	if err != nil {
		return nil, fmt.Errorf("error querying disponibilidad: %w", err)
	}
	defer rows.Close()

	var disponibilidades []domain.DisponibilidadFecha
	for rows.Next() {
		var d domain.DisponibilidadFecha
		err := rows.Scan(&d.Fecha, &d.Disponible, &d.Habitaciones)
		if err != nil {
			return nil, fmt.Errorf("error scanning disponibilidad: %w", err)
		}
		disponibilidades = append(disponibilidades, d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating disponibilidad rows: %w", err)
	}

	return disponibilidades, nil
}

// GetFechasBloqueadas implementa domain.HabitacionRepository
func (r *habitacionRepository) GetFechasBloqueadas(desde, hasta time.Time) (*domain.FechasBloqueadas, error) {
	query := `
		WITH RECURSIVE fechas AS (
			SELECT cast($1 as date) as fecha
			UNION ALL
			SELECT (fecha + interval '1 day')::date
			FROM fechas
			WHERE fecha < cast($2 as date)
		),
		habitaciones_totales AS (
			SELECT COUNT(*) as total
			FROM habitacion h
			WHERE h.estado = 'Disponible'
		),
		habitaciones_ocupadas AS (
			SELECT date(f.fecha) as fecha, 
				   COUNT(DISTINCT rh.habitacionid) as habitaciones_ocupadas
			FROM fechas f
			LEFT JOIN reservaxhabitacion rh ON 
				f.fecha BETWEEN cast(rh.fechaentrada as date) AND cast(rh.fechasalida as date)
			LEFT JOIN reserva r ON r.reservaid = rh.reservaid
				AND rh.estado = 1
				AND r.estado = 'Confirmada'
			GROUP BY f.fecha
			HAVING COUNT(DISTINCT rh.habitacionid) >= (SELECT total FROM habitaciones_totales)
		)
		SELECT fecha::date
		FROM habitaciones_ocupadas
		ORDER BY fecha;`

	rows, err := r.db.Query(query, desde, hasta)
	if err != nil {
		return nil, fmt.Errorf("error querying fechas bloqueadas: %w", err)
	}
	defer rows.Close()

	fechasBloqueadas := &domain.FechasBloqueadas{
		FechasNoDisponibles: make([]time.Time, 0),
	}

	for rows.Next() {
		var fecha time.Time
		if err := rows.Scan(&fecha); err != nil {
			return nil, fmt.Errorf("error scanning fecha bloqueada: %w", err)
		}
		fechasBloqueadas.FechasNoDisponibles = append(fechasBloqueadas.FechasNoDisponibles, fecha)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating fechas bloqueadas: %w", err)
	}

	return fechasBloqueadas, nil
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
			t.precio
		FROM 
			habitacion h
		INNER JOIN 
			tipohabitacion t ON h.tipohabitacionid = t.tipohabitacionid
		WHERE 
			h.estado = 'Disponible'
			AND NOT EXISTS (
				SELECT 1 FROM reservaxhabitacion rh
				JOIN reserva r ON r.reservaid = rh.reservaid
				WHERE rh.habitacionid = h.habitacionid
				AND rh.estado = 1
				AND r.estado = 'Confirmada'
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

// GetRoomTypes returns all available room types
func (r *habitacionRepository) GetRoomTypes() ([]domain.TipoHabitacion, error) {
	query := `
		SELECT 
			t.tipohabitacionid,
			t.titulo,
			t.descripcion,
			t.capacidadadultos,
			t.capacidadninhos,
			t.cantidadcamas,
			t.precio
		FROM tipohabitacion t
		ORDER BY t.tipohabitacionid;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying room types: %w", err)
	}
	defer rows.Close()

	var tipos []domain.TipoHabitacion
	for rows.Next() {
		var t domain.TipoHabitacion
		if err := rows.Scan(
			&t.ID,
			&t.Titulo,
			&t.Descripcion,
			&t.CapacidadAdultos,
			&t.CapacidadNinhos,
			&t.CantidadCamas,
			&t.Precio,
		); err != nil {
			return nil, fmt.Errorf("error scanning room type: %w", err)
		}
		tipos = append(tipos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating room types: %w", err)
	}

	return tipos, nil
}
