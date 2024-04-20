package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"khiemle.dev/golang-api-template/api/routes"
	_authHandler "khiemle.dev/golang-api-template/internal/auth/handler"
	_authService "khiemle.dev/golang-api-template/internal/auth/service"
	_todoHandler "khiemle.dev/golang-api-template/internal/todo/handler"
	_todoService "khiemle.dev/golang-api-template/internal/todo/service"
	_userService "khiemle.dev/golang-api-template/internal/user/service"
	util "khiemle.dev/golang-api-template/pkg/util"
	"khiemle.dev/golang-api-template/pkg/util/token"
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

	// Token
	tokenMaker := token.NewTokenMaker(s.cfg)

	// Services
	todoService := _todoService.NewTodoService(s.db)
	userService := _userService.NewUserService(s.db)
	loginSessionService := _authService.NewLoginSessionService(s.db)
	authService := _authService.NewAuthService(s.db, s.cfg, userService, loginSessionService, tokenMaker)

	// Handlers
	todoHandler := _todoHandler.NewTodoHandler(todoService)
	authHandler := _authHandler.NewAuthHandler(s.cfg, authService, loginSessionService)

	// Routes for v1 endpoints
	v1 := s.router.Group("/v1")
	{
		// Setup todoGroupRouter
		todoGroup := v1.Group("/todos")
		routes.SetupTodoRouter(todoGroup, todoHandler)

		// Setup authGroupRouter
		authGroup := v1.Group("/auth")
		routes.SetupAuthRouter(authGroup, authHandler, tokenMaker, loginSessionService, userService)
	}

	log.Info().Msg("Routes setup complete!")
}
