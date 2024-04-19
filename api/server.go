package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/api/routes"
	"khiemle.dev/golang-api-template/internal/todo/handler"
	"khiemle.dev/golang-api-template/internal/todo/service"
	util "khiemle.dev/golang-api-template/pkg/util"
)

// Server represents the HTTP server.
type Server struct {
	router *gin.Engine
	cfg    *util.Config
	db     *gorm.DB
}

// NewServer creates a new HTTP server instance.
func NewServer() *Server {
	return &Server{}
}

// Initialize initializes the HTTP server.
func (s *Server) Initialize(cfg *util.Config, db *gorm.DB) error {
	s.cfg = cfg
	s.db = db

	// Create a new Gin router
	s.router = gin.Default()

	// Setup routes
	s.setupRoutes()

	return nil
}

func (s *Server) StartServer() error {
	return s.router.Run(s.cfg.HTTPServerAddress)
}

// Setup routes for the server.
func (s *Server) setupRoutes() {
	log.Info().Msg("Setting up routes...")

	// Routes for health check
	s.router.GET("/health-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Services
	todoService := service.NewTodoService(s.db)

	// Handlers
	todoHandler := handler.NewTodoHandler(todoService)

	// Routes for v1 endpoints
	v1 := s.router.Group("/v1")
	{
		// Setup todoGroupRouter
		todoGroup := v1.Group("/todos")
		routes.SetupTodoRouter(todoGroup, todoHandler)
	}

	log.Info().Msg("Routes setup complete!")
}
