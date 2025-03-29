package routers

import (
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/freterapido-backend-api/infrastructure/logger"
	"github.com/thalesmacedo1/freterapido-backend-api/interfaces/api/handlers"
	"github.com/thalesmacedo1/freterapido-backend-api/interfaces/middleware"
)

func Router(quoteHandler *handlers.QuoteHandler, metricsHandler *handlers.MetricsHandler, logger logger.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware(logger))

	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Quote Endpoints

	router.GET("/quote", quoteHandler.MakeQuotation)
	router.POST("/metrics", metricsHandler.GetMetrics)

	return router
}
