package routes

import (
	"github.com/gin-gonic/gin"
	_authHandler "khiemle.dev/golang-api-template/internal/auth/handler"
	"khiemle.dev/golang-api-template/internal/auth/service"
	_todoHandler "khiemle.dev/golang-api-template/internal/todo/handler"
	_userService "khiemle.dev/golang-api-template/internal/user/service"
	"khiemle.dev/golang-api-template/pkg/middleware"
	"khiemle.dev/golang-api-template/pkg/util/token"
)

func SetupTodoRouter(todoGroup *gin.RouterGroup, todoHandler _todoHandler.TodoHandler) {
	todoGroup.GET("/", todoHandler.ListTodoHandler)
	todoGroup.GET("/:id", todoHandler.GetByIdHandler)
	todoGroup.POST("/", todoHandler.CreateTodoHandler)
	todoGroup.PATCH("/:id", todoHandler.UpdateTodoHandler)
	todoGroup.DELETE("/:id", todoHandler.DeleteTodoHandler)
}

func SetupAuthRouter(
	authGroup *gin.RouterGroup,
	authHandler _authHandler.AuthHandler,
	tokenMaker token.TokenMaker,
	loginSessionService service.LoginSessionService,
	userService _userService.UserService,
) {
	authMiddleware := middleware.AuthorizationMiddleware(tokenMaker, loginSessionService, userService)

	authGroup.POST("/login", authHandler.LoginHandler)
	authGroup.POST("/register", authHandler.RegisterHandler)
	authGroup.DELETE("/logout", authMiddleware, authHandler.LogoutHandler)
	authGroup.GET("/verify_access_token", authMiddleware, authHandler.VerifyAccessToken)
}
