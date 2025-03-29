package main

import (
	"log"

	"github.com/thalesmacedo1/freterapido-backend-api/application/usecases"
	"github.com/thalesmacedo1/freterapido-backend-api/config"
	"github.com/thalesmacedo1/freterapido-backend-api/infrastructure/database/postgres/repositories"
	"github.com/thalesmacedo1/freterapido-backend-api/infrastructure/logger"
	"github.com/thalesmacedo1/freterapido-backend-api/interfaces/api/handlers"
	"github.com/thalesmacedo1/freterapido-backend-api/interfaces/routers"

	"github.com/thalesmacedo1/freterapido-backend-api/infrastructure/database/postgres"
)

func main() {
	// Carrega as configurações
	if err := config.LoadConfig(".env.example"); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Inicializa o logger
	logr := logger.NewLogrusLogger()

	// Inicializa o cliente Postgres
	postgresClient, err := postgres.NewPostgresClient(config.Settings.Neo4jURI, config.Settings.Neo4jUser, config.Settings.Neo4jPassword)
	if err != nil {
		logr.Fatalf("Failed to initialize Postgres client: %v", err)
	}
	defer postgresClient.Close()

	// Inicializa os repositórios
	metricsRepo := repositories.NewPostgresMetricsRepository(postgresClient.Driver)
	quoteRepo := repositories.NewPostgresQuoteRepository(postgresClient.Driver)

	// // Inicializa o cliente Redis
	// redisClient, err := redis.NewRedisClient(config.Settings.RedisHost, config.Settings.RedisPassword, config.Settings.RedisDB)
	// if err != nil {
	// 	logr.Fatalf("Failed to initialize Redis client: %v", err)
	// }
	// defer redisClient.Close()

	// Inicializa os use cases
	getShippingQuotationUC := usecases.NewGetShippingQuotationUseCase(quoteRepo)
	getLastQuotesUC := usecases.NewGetLastQuotesUseCase(metricsRepo)

	// Inicializa os handlers
	quoteHandler := handlers.NewQuoteHandler(getShippingQuotationUC, logr)
	metricsHandler := handlers.NewMetricsHandler(getLastQuotesUC, logr)

	// Configura o roteador usando Gin
	router := routers.Router(quoteHandler, metricsHandler, logr)

	// Inicia o servidor HTTP
	logr.Info("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		logr.Fatalf("Server failed to start: %v", err)
	}
}
