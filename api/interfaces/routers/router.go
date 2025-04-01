package routers

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thalesmacedo1/freterapido-backend-api/api/application/usecases"
	"github.com/thalesmacedo1/freterapido-backend-api/api/interfaces/api"
)

// SetupRouter configures the API routes
func SetupRouter(
	getShippingQuotationUseCase *usecases.GetShippingQuotationUseCase,
	getMetricsUseCase *usecases.GetMetricsUseCase,
) *gin.Engine {
	router := gin.Default()

	// Create controllers
	quoteController := api.NewQuoteController(getShippingQuotationUseCase)
	metricsController := api.NewMetricsController(getMetricsUseCase)

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// API routes group
	apiGroup := router.Group("/")
	{
		// Quote route
		apiGroup.POST("/quote", quoteController.GetQuote)

		// Metrics route
		apiGroup.GET("/metrics", metricsController.GetMetrics)
	}

	return router
}
