package main

import (
	"database/sql"
	"log"

	"github.com/Maxito7/hotel_backend/internal/application"
	"github.com/Maxito7/hotel_backend/internal/config"
	"github.com/Maxito7/hotel_backend/internal/infrastructure/repository"
	handlers "github.com/Maxito7/hotel_backend/internal/interfaces/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq" // Driver de PostgreSQL
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Conectar a la base de datos
	db, err := sql.Open("postgres", cfg.GetDBConnString())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Verificar conexión a la base de datos
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Inicializar Fiber
	app := fiber.New()

	// Configurar CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // URL de tu frontend Next.js
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400, // 24 horas en segundos
	}))

	// Inicializar dependencias
	habitacionRepo := repository.NewHabitacionRepository(db)
	habitacionService := application.NewHabitacionService(habitacionRepo)
	habitacionHandler := handlers.NewHabitacionHandler(habitacionService)

	contactRepo := repository.NewContactRepository(db)
	contactService := application.NewContactService(contactRepo)
	contactHandler := handlers.NewContactHandler(contactService)

	// Configurar rutas
	api := app.Group("/api")
	habitaciones := api.Group("/habitaciones")

	// Rutas de habitaciones
	habitaciones.Get("/", habitacionHandler.GetAllRooms)
	habitaciones.Get("/disponibles", habitacionHandler.GetAvailableRooms)
	habitaciones.Get("/fechas-bloqueadas", habitacionHandler.GetFechasBloqueadas)

	// Rutas de contacto
	contacto := api.Group("/contact")
	contacto.Post("/", contactHandler.Create)
	contacto.Get("/", contactHandler.List)
	contacto.Patch("/:id/estado", contactHandler.UpdateEstado)

	// Iniciar servidor
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
