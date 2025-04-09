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

	_ "github.com/thalesmacedo1/freterapido-backend-api/docs"
)

// @title API de Cotação de Frete
// @version 1.0
// @description API para consulta de valores de frete através de integrações com transportadoras.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.freterapido.com
// @contact.email suporte@freterapido.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	dbHost := getEnv("POSTGRES_HOST")
	dbPort := getEnv("POSTGRES_PORT")
	dbUser := getEnv("POSTGRES_USER")
	dbPassword := getEnv("POSTGRES_PASSWORD")
	dbName := getEnv("POSTGRES_DB")

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

	router := routers.SetupRouter(getShippingQuotationUseCase, getMetricsUseCase)

	port := getEnv("PORT", "3000")

	log.Printf("Starting server on port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
