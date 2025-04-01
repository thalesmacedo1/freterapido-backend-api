package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thalesmacedo1/freterapido-backend-api/api/application/usecases"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
	"github.com/thalesmacedo1/freterapido-backend-api/api/infrastructure/database"
	"github.com/thalesmacedo1/freterapido-backend-api/api/interfaces/routers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Obtém configurações do banco de dados
	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbUser := getEnv("POSTGRES_USER", "postgres")
	dbPassword := getEnv("POSTGRES_PASSWORD", "postgres")
	dbName := getEnv("POSTGRES_DB", "freterapido")

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	log.Printf("Connecting to database: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(&domain.QuoteResponse{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories
	quoteRepository := database.NewQuoteRepository(db)
	metricsRepository := database.NewMetricsRepository(db)

	// Create use cases
	getShippingQuotationUseCase := usecases.NewGetShippingQuotationUseCase(quoteRepository)
	getMetricsUseCase := usecases.NewGetMetricsUseCase(metricsRepository)

	// Setup router
	router := routers.SetupRouter(getShippingQuotationUseCase, getMetricsUseCase)

	// Get server port from environment or use default
	port := getEnv("PORT", "8080")

	log.Printf("Starting server on port %s", port)

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
