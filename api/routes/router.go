package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	todoHandler "khiemle.dev/golang-api-template/internal/todo/handler"
	"khiemle.dev/golang-api-template/pkg/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	log.Info().Msg("Setting up routes...")

	// Routes for health check
	r.GET("/health-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Routes for v1 endpoints
	v1 := r.Group("/v1")
	{
		todoGroup := v1.Group("/protected_items")

		// Setup auth middleware
		todoGroup.Use(middleware.AuthorizationMiddleware())

		todoGroup.GET("/list", todoHandler.ListAllTodos)
	}

	log.Info().Msg("Routes setup complete!")

	return r
}
