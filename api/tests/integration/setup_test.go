package integration

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/freterapido-backend-api/api/application/usecases"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
	"github.com/thalesmacedo1/freterapido-backend-api/api/infrastructure/database"
	"github.com/thalesmacedo1/freterapido-backend-api/api/interfaces/routers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testDB                *gorm.DB
	testQuoteRepository   domain.QuoteRepository
	testMetricsRepository domain.MetricsRepository
	testRouter            *gin.Engine
)

// setupTestDB establishes a test database connection
func setupTestDB() (*gorm.DB, error) {
	// Get connection details from environment variables
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=freterapido_test port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	if err := db.AutoMigrate(&domain.QuoteResponse{}); err != nil {
		return nil, err
	}

	return db, nil
}

// cleanupDB clears all test data
func cleanupDB(db *gorm.DB) error {
	return db.Exec("TRUNCATE TABLE quote_responses CASCADE").Error
}

// setupTestEnvironment initializes the test environment
func setupTestEnvironment() error {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup database
	var err error
	testDB, err = setupTestDB()
	if err != nil {
		return err
	}

	// Initialize repositories
	testQuoteRepository = database.NewQuoteRepository(testDB)
	testMetricsRepository = database.NewMetricsRepository(testDB)

	// Initialize use cases
	getShippingQuotationUseCase := usecases.NewGetShippingQuotationUseCase(testQuoteRepository)
	getMetricsUseCase := usecases.NewGetMetricsUseCase(testMetricsRepository)

	// Setup router
	testRouter = routers.SetupRouter(getShippingQuotationUseCase, getMetricsUseCase)

	return nil
}

// TestMain is the entry point for all tests in this package
func TestMain(m *testing.M) {
	// Skip integration tests if not specifically enabled
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		log.Println("Skipping integration tests. Set RUN_INTEGRATION_TESTS=true to run them.")
		os.Exit(0)
	}

	// Setup test environment
	if err := setupTestEnvironment(); err != nil {
		log.Fatalf("Failed to setup test environment: %v", err)
	}

	// Clean up before tests
	if err := cleanupDB(testDB); err != nil {
		log.Fatalf("Failed to clean database before tests: %v", err)
	}

	// Run tests
	exitCode := m.Run()

	// Clean up after tests
	if err := cleanupDB(testDB); err != nil {
		log.Printf("Failed to clean database after tests: %v", err)
	}

	os.Exit(exitCode)
}
