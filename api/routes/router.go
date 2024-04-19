package routes

import (
	"github.com/gin-gonic/gin"
	_authHandler "khiemle.dev/golang-api-template/internal/auth/handler"
	_todoHandler "khiemle.dev/golang-api-template/internal/todo/handler"
)

func SetupTodoRouter(todoGroup *gin.RouterGroup, todoHandler _todoHandler.TodoHandler) {
	todoGroup.GET("/", todoHandler.ListTodoHandler)
	todoGroup.GET("/:id", todoHandler.GetByIdHandler)
	todoGroup.POST("/", todoHandler.CreateTodoHandler)
	todoGroup.PATCH("/:id", todoHandler.UpdateTodoHandler)
	todoGroup.DELETE("/:id", todoHandler.DeleteTodoHandler)
}

func SetupAuthRouter(authGroup *gin.RouterGroup, authHandler _authHandler.AuthHandler) {
	authGroup.POST("/login", authHandler.LoginHandler)
	authGroup.POST("/register", authHandler.RegisterHandler)
}
