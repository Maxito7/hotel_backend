package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	ServerPort   string
	TavilyAPIKey string
	OpenAIAPIKey string
}

func LoadConfig() (*Config, error) {
	// Cargar variables de entorno desde el archivo .env
	// No retornar error si el archivo no existe, permitir variables de entorno del sistema
	_ = godotenv.Load()

	config := &Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", ""),
		DBName:       getEnv("DB_NAME", "postgres"),
		ServerPort:   getEnv("SERVER_PORT", "8000"),
		TavilyAPIKey: getEnv("TAVILY_API_KEY", ""),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
	}

	// Validar que las variables requeridas no estén vacías
	if config.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	return config, nil
}

func (c *Config) GetDBConnString() string {
	// Asegurar conexión SSL y timeout
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require connect_timeout=5",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

// String implementa la interfaz Stringer para evitar que se impriman datos sensibles en logs
func (c Config) String() string {
	return fmt.Sprintf("Config{DBHost: %s, DBPort: %s, DBUser: %s, DBPassword: [HIDDEN], DBName: %s, ServerPort: %s}",
		c.DBHost, c.DBPort, c.DBUser, c.DBName, c.ServerPort)
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
